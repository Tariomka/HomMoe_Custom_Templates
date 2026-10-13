package editor_state_model

import (
	"reflect"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template_entity"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/linq"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
)

type ManualEditSettings struct {
	ManualZones       []template_model.Zone
	ManualConnections []ManualConnectionSave
}

func (this ManualEditSettings) Equals(other ManualEditSettings) bool {
	return manualZonesEqual(this.ManualZones, other.ManualZones) &&
		manualConnectionsEqual(this.ManualConnections, other.ManualConnections)
}

func ToManualEditSettingsModel(entity editor_state.ManualEditSettings) ManualEditSettings {
	return ManualEditSettings{
		ManualZones:       linq.SelectSlice(entity.Zones, zoneFromAllOverridden),
		ManualConnections: linq.SelectSlice(entity.Connections, connectionFromAllOverridden),
	}
}

func ToManualEditSettingsEntity(model ManualEditSettings) editor_state.ManualEditSettings {
	return editor_state.ManualEditSettings{
		Zones:       linq.SelectSlice(model.ManualZones, zoneToAllOverridden),
		Connections: linq.SelectSlice(model.ManualConnections, connectionToAllOverridden),
	}
}

func zoneToAllOverridden(zone template_model.Zone) editor_state.ManualZoneEdit {
	var emptyBase template_entity.Zone
	overrides, unset := DiffFieldOverrides[template_entity.Zone, editor_state.ZoneFieldOverrides](
		emptyBase, template_model.ToZoneEntity(zone))
	return editor_state.ManualZoneEdit{
		Name:              zone.Name,
		Quality:           toQualityOrdinal(zone.Quality),
		GeneratorPosition: helpers.ClonePointer(zone.GeneratorPosition),
		GeneratorRing:     helpers.ClonePointer(zone.GeneratorRing),
		ManualPosition:    helpers.ClonePointer(zone.ManualPosition),
		FieldOverrides:    overrides,
		Unset:             unset,
	}
}

func zoneFromAllOverridden(edit editor_state.ManualZoneEdit) template_model.Zone {
	var base template_entity.Zone
	base.Name = edit.Name
	zone := template_model.ToZoneModel(ApplyFieldOverrides(base, edit.FieldOverrides, edit.Unset))
	zone.GeneratorPosition = helpers.ClonePointer(edit.GeneratorPosition)
	zone.GeneratorRing = helpers.ClonePointer(edit.GeneratorRing)
	zone.ManualPosition = helpers.ClonePointer(edit.ManualPosition)
	zone.Quality = fromQualityOrdinal(edit.Quality)
	return zone
}

func connectionToAllOverridden(save ManualConnectionSave) editor_state.ManualConnectionEdit {
	var emptyBase template_entity.Connection
	overrides, unset := DiffFieldOverrides[template_entity.Connection, editor_state.ConnectionFieldOverrides](
		emptyBase, cloneConnection(save.Connection))
	return editor_state.ManualConnectionEdit{
		Name:           save.Connection.Name,
		From:           save.Connection.From,
		To:             save.Connection.To,
		ConnectionType: save.Connection.ConnectionType,
		IsUserAdded:    save.IsUserAdded,
		FieldOverrides: overrides,
		Unset:          unset,
	}
}

func connectionFromAllOverridden(edit editor_state.ManualConnectionEdit) ManualConnectionSave {
	var base template_entity.Connection
	base.Name, base.From, base.To, base.ConnectionType = edit.Name, edit.From, edit.To, edit.ConnectionType
	return ManualConnectionSave{
		Connection:  cloneConnection(ApplyFieldOverrides(base, edit.FieldOverrides, edit.Unset)),
		IsUserAdded: edit.IsUserAdded}
}

func manualZonesEqual(left, right []template_model.Zone) bool {
	if len(left) == 0 || len(right) == 0 {
		return len(left) == len(right)
	}

	return reflect.DeepEqual(ToManualZoneSaveEntities(left), ToManualZoneSaveEntities(right))
}

func manualConnectionsEqual(left, right []ManualConnectionSave) bool {
	if len(left) == 0 || len(right) == 0 {
		return len(left) == len(right)
	}

	return reflect.DeepEqual(left, right)
}
