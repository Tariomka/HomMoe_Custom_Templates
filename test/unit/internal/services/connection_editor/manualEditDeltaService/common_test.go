// Package manualEditDeltaService_test contains shared arrangement helpers for
// the manualEditDeltaService.go unit tests.
package manualEditDeltaService_test

import (
	"reflect"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/connection_editor"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/connection_editor/connection_editor_interfaces"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/generation_tuning"
	zone_services "github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
)

const (
	spawnZoneName    = "Spawn-A"
	goldZoneName     = "Neutral-G"
	silverZoneName   = "Neutral-H"
	spawnToGoldName  = "A-G"
	goldToSilverName = "G-H"
	addedZoneLabel   = "Q"
	addedZoneName    = "Neutral-Q"
)

// newService builds the service with the same collaborators the application wires.
func newService() connection_editor_interfaces.IManualEditDeltaService {
	zoneEditor := test_helpers.NewZoneEditorService()
	tierService := zone_services.NewZoneTierService()
	tuningFactory := generation_tuning.NewGenerationTuningFactory()
	return connection_editor.NewManualEditDeltaService(
		zoneEditor,
		connection_editor.NewConnectionEditorService(tierService),
		connection_editor.NewManualReapplyService(
			zoneEditor, zone_services.NewCastleFactory(), tierService, tuningFactory),
		tierService,
		tuningFactory,
	)
}

func newConfiguration() *config.GeneratorConfig {
	return config.NewGeneratorConfig()
}

func newTuning(configuration *config.GeneratorConfig, zoneCount int) models.GenerationTuning {
	return generation_tuning.NewGenerationTuningFactory().Create(configuration, zoneCount)
}

// factoryZone builds a zone exactly as the service derives an added or
// re-profiled one.
func factoryZone(
	label string,
	quality neutral_zone.Quality,
	castles int,
	configuration *config.GeneratorConfig,
	zoneCount int) template_model.Zone {
	return test_helpers.NewZoneEditorService().NewDefaultNeutralZone(
		label, quality, castles, configuration.GenerateRoads, newTuning(configuration, zoneCount))
}

// freshLayout is a small generated layout: a spawn, a Gold and a Silver
// neutral zone, all stamped with generator positions, and two connections.
func freshLayout(configuration *config.GeneratorConfig) ([]template_model.Zone, []template_model.Connection) {
	spawn := test_helpers.NewZoneEditorService().NewDefaultNeutralZone(
		"A", neutral_zone.QualityLow, 1, configuration.GenerateRoads, newTuning(configuration, 3))
	spawn.Name = spawnZoneName
	spawn.Quality = nil
	spawn.MainObjects = []template_model.MainObject{{Type: "Spawn", Spawn: "Player1"}}
	gold := factoryZone("G", neutral_zone.QualityHigh, 1, configuration, 3)
	gold.Size = 1.25
	silver := factoryZone("H", neutral_zone.QualityMedium, 0, configuration, 3)
	silver.Size = 0.75
	zones := []template_model.Zone{spawn, gold, silver}
	for index := range zones {
		zones[index].GeneratorPosition = new(data.NewVec2(0.2*float64(index+1), 0.5))
		zones[index].GeneratorRing = new(index)
	}

	connections := []template_model.Connection{
		{Name: spawnToGoldName, From: spawnZoneName, To: goldZoneName, ConnectionType: "Direct", GuardValue: 1000},
		{Name: goldToSilverName, From: goldZoneName, To: silverZoneName, ConnectionType: "Direct", GuardValue: 2000},
	}
	return zones, connections
}

// captureRequest records edited zones and connections against freshLayout.
func captureRequest(
	configuration *config.GeneratorConfig,
	edit func(zones []template_model.Zone, connections []template_model.Connection) (
		[]template_model.Zone, []template_model.Connection)) models.ManualEditCaptureRequest {
	freshZones, freshConnections := freshLayout(configuration)
	zones, connections := freshLayout(configuration)
	zones, connections = edit(zones, connections)
	return models.ManualEditCaptureRequest{
		FreshZones:       freshZones,
		FreshConnections: freshConnections,
		Zones:            zones,
		Connections:      connections,
		Configuration:    configuration,
	}
}

func unchanged(
	zones []template_model.Zone,
	connections []template_model.Connection) ([]template_model.Zone, []template_model.Connection) {
	return zones, connections
}

// expandRequest rebuilds stored edits on freshLayout.
func expandRequest(
	configuration *config.GeneratorConfig,
	edits editor_state_model.ManualEditSettings) models.ManualEditExpandRequest {
	freshZones, freshConnections := freshLayout(configuration)
	return models.ManualEditExpandRequest{
		FreshZones:       freshZones,
		FreshConnections: freshConnections,
		ManualEdits:      edits,
		Configuration:    configuration,
	}
}

func zoneEdit(settings editor_state_model.ManualEditSettings, name string) editor_state_model.ManualZoneEdit {
	for _, edit := range settings.ManualZones {
		if edit.Name == name {
			return edit
		}
	}
	return editor_state_model.ManualZoneEdit{}
}

func zoneNamed(zones []template_model.Zone, name string) template_model.Zone {
	for _, zone := range zones {
		if zone.Name == name {
			return zone
		}
	}
	return template_model.Zone{}
}

// overriddenZoneNames lists the zones whose edit stores any override or unset field.
func overriddenZoneNames(settings editor_state_model.ManualEditSettings) []string {
	var names []string
	for _, edit := range settings.ManualZones {
		if !reflect.ValueOf(edit.FieldOverrides).IsZero() || len(edit.Unset) > 0 {
			names = append(names, edit.Name)
		}
	}
	return names
}

// overriddenConnectionNames lists the connections whose edit stores any override or unset field.
func overriddenConnectionNames(settings editor_state_model.ManualEditSettings) []string {
	var names []string
	for _, edit := range settings.ManualConnections {
		if !reflect.ValueOf(edit.FieldOverrides).IsZero() || len(edit.Unset) > 0 {
			names = append(names, edit.Name)
		}
	}
	return names
}

// userAddedConnection builds a connection from the spawn the way the editor adds one.
func userAddedConnection(
	to string,
	zones []template_model.Zone,
	configuration *config.GeneratorConfig) template_model.Connection {
	playerZoneNames := map[string]bool{spawnZoneName: true}
	connection := connection_editor.NewConnectionEditorService(zone_services.NewZoneTierService()).
		NewDefaultConnection(spawnZoneName, to, zones, playerZoneNames)
	test_helpers.NewZoneEditorService().ApplyConnectionRoadPolicy(&connection, configuration.GenerateRoads)
	return connection
}
