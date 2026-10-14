package connection_editor

import (
	"fmt"
	"slices"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/zone_helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/preview"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/connection_editor/connection_editor_interfaces"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/generation_tuning"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones/zone_interfaces"
)

// missingZoneCastles is the castle count a zone the generator no longer
// produces is rebuilt with.
const missingZoneCastles = 1

// ManualEditDeltaService derives every manual zone from the same-named zone of
// a fresh generation, or from the editor's neutral-zone factory when the zone
// was added or re-profiled, and every connection from the same-named fresh
// connection with the same endpoints, or from the editor's connection factory.
// Only the differences from that derivation are stored.
type ManualEditDeltaService struct {
	zoneEditor       connection_editor_interfaces.IZoneEditorService
	connectionEditor connection_editor_interfaces.IConnectionEditorService
	manualReapply    connection_editor_interfaces.IManualReapplyService
	tierService      zone_interfaces.IZoneTierService
	tuningFactory    generation_tuning.IGenerationTuningFactory
}

func NewManualEditDeltaService(
	zoneEditor connection_editor_interfaces.IZoneEditorService,
	connectionEditor connection_editor_interfaces.IConnectionEditorService,
	manualReapply connection_editor_interfaces.IManualReapplyService,
	tierService zone_interfaces.IZoneTierService,
	tuningFactory generation_tuning.IGenerationTuningFactory) connection_editor_interfaces.IManualEditDeltaService {
	return &ManualEditDeltaService{
		zoneEditor:       zoneEditor,
		connectionEditor: connectionEditor,
		manualReapply:    manualReapply,
		tierService:      tierService,
		tuningFactory:    tuningFactory,
	}
}

func (this *ManualEditDeltaService) Capture(
	request models.ManualEditCaptureRequest) editor_state_model.ManualEditSettings {
	freshZones := indexByName(request.FreshZones, func(zone template_model.Zone) string { return zone.Name })
	legacyUnstamped := !slices.ContainsFunc(request.Zones, func(zone template_model.Zone) bool {
		return zone.GeneratorPosition != nil
	})
	tuning := this.tuningFactory.Create(request.Configuration, len(request.Zones))

	var edits editor_state_model.ManualEditSettings
	for _, zone := range request.Zones {
		freshZone, found := freshZones[zone.Name]
		added := !found || (!legacyUnstamped && zone.GeneratorPosition == nil)
		if !added && !this.isReprofiled(zone, freshZone) {
			edits.ManualZones = append(edits.ManualZones,
				editor_state_model.NewManualZoneEdit(freshZone, zone, request.NonDeterministic))
			continue
		}

		castles := this.zoneEditor.CountZoneCastles(zone)
		derived := this.factoryZone(zone.Name, zone.Quality, castles, tuning, request.Configuration.GenerateRoads)
		edit := editor_state_model.NewManualZoneEdit(derived, zone, request.NonDeterministic)
		if added {
			edit = edit.AsAdded(castles)
		} else {
			edit = edit.AsReprofiled(castles)
		}
		edits.ManualZones = append(edits.ManualZones, edit)
	}

	deriveConnection := this.newConnectionSource(
		request.FreshConnections,
		request.Zones,
		request.Configuration.GenerateRoads,
		request.NonDeterministic,
	)
	for _, connection := range request.Connections {
		edits.ManualConnections = append(edits.ManualConnections,
			editor_state_model.NewManualConnectionEdit(deriveConnection(connection), connection))
	}
	return edits
}

func (this *ManualEditDeltaService) Expand(request models.ManualEditExpandRequest) models.ManualEditExpansion {
	freshZones := indexByName(request.FreshZones, func(zone template_model.Zone) string { return zone.Name })
	tuning := this.tuningFactory.Create(request.Configuration, len(request.ManualEdits.ManualZones))
	edits := request.ManualEdits.Clone()

	zones, eligible, warnings := this.expandZones(edits.ManualZones, freshZones, tuning, request)
	if request.CastleChanges.Any() {
		this.manualReapply.ApplyCastleSettingChanges(zones, eligible, request.CastleChanges, request.Configuration)
		for index, zone := range zones {
			if eligible[index] {
				edits.ManualZones[index] = this.recaptureZone(
					edits.ManualZones[index],
					zone,
					freshZones,
					tuning,
					request,
				)
			}
		}
	}

	deriveConnection := this.newConnectionSource(
		request.FreshConnections,
		zones,
		request.Configuration.GenerateRoads,
		request.NonDeterministic,
	)
	connections := make([]template_model.Connection, len(edits.ManualConnections))
	for index, edit := range edits.ManualConnections {
		connections[index] = edit.Apply(deriveConnection(edit.Structure()))
	}

	return models.ManualEditExpansion{
		Zones:       zones,
		Connections: connections,
		ManualEdits: edits,
		Changed:     !edits.Equals(request.ManualEdits),
		Warnings:    warnings,
	}
}

// expandZones rebuilds every zone on its derivation and reports which own
// their castles (eligible for the castle reapply). A generated zone the fresh
// generation no longer holds is rebuilt by the factory and re-recorded in
// edits as added.
func (this *ManualEditDeltaService) expandZones(
	edits []editor_state_model.ManualZoneEdit,
	freshZones map[string]template_model.Zone,
	tuning models.GenerationTuning,
	request models.ManualEditExpandRequest) ([]template_model.Zone, []bool, []string) {
	generateRoads := request.Configuration.GenerateRoads
	var warnings []string
	zones := make([]template_model.Zone, len(edits))
	eligible := make([]bool, len(edits))
	for index, edit := range edits {
		if castles, ok := edit.GetCastles(); ok || edit.IsAdded() {
			zones[index] = edit.Apply(this.factoryZone(edit.Name, edit.GetQuality(), castles, tuning, generateRoads))
			eligible[index] = true
			continue
		}

		if freshZone, found := freshZones[edit.Name]; found {
			zones[index] = edit.Apply(freshZone)
			eligible[index] = edit.HasMainObjectsOverride()
			continue
		}

		derived := this.factoryZone(edit.Name, edit.GetQuality(), missingZoneCastles, tuning, generateRoads)
		zones[index] = edit.Apply(derived)
		edits[index] = editor_state_model.NewManualZoneEdit(derived, zones[index], request.NonDeterministic).
			AsAdded(missingZoneCastles)
		eligible[index] = true
		warnings = append(warnings,
			fmt.Sprintf("zone %q is no longer generated; kept as an editor-made zone", edit.Name))
	}
	return zones, eligible, warnings
}

// recaptureZone records a zone again after the castle reapply rebuilt it, so
// the stored castle count and overrides follow the new castles.
func (this *ManualEditDeltaService) recaptureZone(
	edit editor_state_model.ManualZoneEdit,
	zone template_model.Zone,
	freshZones map[string]template_model.Zone,
	tuning models.GenerationTuning,
	request models.ManualEditExpandRequest) editor_state_model.ManualZoneEdit {
	if _, ok := edit.GetCastles(); !ok && !edit.IsAdded() {
		return editor_state_model.NewManualZoneEdit(freshZones[edit.Name], zone, request.NonDeterministic)
	}

	castles := this.zoneEditor.CountZoneCastles(zone)
	derived := this.factoryZone(zone.Name, zone.Quality, castles, tuning, request.Configuration.GenerateRoads)
	recaptured := editor_state_model.NewManualZoneEdit(derived, zone, request.NonDeterministic)
	if edit.IsAdded() {
		return recaptured.AsAdded(castles)
	}
	return recaptured.AsReprofiled(castles)
}

// isReprofiled reports whether a generated neutral zone's tier or castle count
// was changed away from the fresh zone of the same name.
func (this *ManualEditDeltaService) isReprofiled(zone, freshZone template_model.Zone) bool {
	return zone_helpers.GetZoneTypeFromName(zone.Name) == preview.ZoneTypeNeutral &&
		(this.tierService.ResolveQuality(zone) != this.tierService.ResolveQuality(freshZone) ||
			this.zoneEditor.CountZoneCastles(zone) != this.zoneEditor.CountZoneCastles(freshZone))
}

// factoryZone builds the derivation of an added or re-profiled zone. A zone
// without a recorded tier derives from the Medium profile.
func (this *ManualEditDeltaService) factoryZone(
	name string,
	quality *neutral_zone.Quality,
	castles int,
	tuning models.GenerationTuning,
	generateRoads bool) template_model.Zone {
	tier := neutral_zone.QualityMedium
	if quality != nil {
		tier = *quality
	}

	return this.zoneEditor.NewDefaultNeutralZone(helpers.GetZoneLabel(name), tier, castles, generateRoads, tuning)
}

// newConnectionSource returns the connection derivation for one layout. The
// connections of a non-deterministic layout differ between generations, so
// those always derive from the factory.
func (this *ManualEditDeltaService) newConnectionSource(
	freshConnections []template_model.Connection,
	zones []template_model.Zone,
	generateRoads bool,
	nonDeterministic bool) func(structure template_model.Connection) template_model.Connection {
	if nonDeterministic {
		freshConnections = nil
	}
	fresh := indexByName(freshConnections, func(connection template_model.Connection) string {
		return connection.Name
	})
	playerZoneNames := make(map[string]bool)
	for _, zone := range zones {
		if zone_helpers.IsZoneNamePlayer(zone.Name) {
			playerZoneNames[zone.Name] = true
		}
	}

	return func(structure template_model.Connection) template_model.Connection {
		return this.deriveConnection(structure, fresh, zones, playerZoneNames, generateRoads)
	}
}

// deriveConnection derives a connection from its structure: the same-named
// fresh connection with the same endpoints, unless the connection was added
// in the editor; otherwise the editor's connection factory.
func (this *ManualEditDeltaService) deriveConnection(
	structure template_model.Connection,
	fresh map[string]template_model.Connection,
	zones []template_model.Zone,
	playerZoneNames map[string]bool,
	generateRoads bool) template_model.Connection {
	if freshConnection, found := fresh[structure.Name]; found && !structure.IsUserAdded &&
		freshConnection.From == structure.From && freshConnection.To == structure.To {
		return freshConnection
	}

	derived := this.connectionEditor.NewDefaultConnection(structure.From, structure.To, zones, playerZoneNames)
	this.zoneEditor.ApplyConnectionRoadPolicy(&derived, generateRoads)
	if derived.ConnectionType != structure.ConnectionType {
		this.zoneEditor.ChangeConnectionType(&derived, structure.ConnectionType, generateRoads)
	}
	return derived
}

func indexByName[Item any](items []Item, name func(Item) string) map[string]Item {
	indexed := make(map[string]Item, len(items))
	for _, item := range items {
		if _, exists := indexed[name(item)]; !exists {
			indexed[name(item)] = item
		}
	}
	return indexed
}
