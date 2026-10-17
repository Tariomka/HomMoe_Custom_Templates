package editor_state_model

import (
	"reflect"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/linq"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
)

type ManualEditSettings struct {
	ManualZones       []ManualZoneEdit
	ManualConnections []ManualConnectionEdit
}

func ManualEditSettingsFromSnapshot(
	zones []template_model.Zone,
	connections []template_model.Connection) ManualEditSettings {
	var emptyZone template_model.Zone
	var emptyConnection template_model.Connection
	return ManualEditSettings{
		ManualZones: linq.SelectSlice(zones, func(zone template_model.Zone) ManualZoneEdit {
			return NewManualZoneEdit(emptyZone, zone, false)
		}),
		ManualConnections: linq.SelectSlice(
			connections,
			func(connection template_model.Connection) ManualConnectionEdit {
				return NewManualConnectionEdit(emptyConnection, connection)
			},
		),
	}
}

func (this ManualEditSettings) ToSnapshot() ([]template_model.Zone, []template_model.Connection) {
	var emptyZone template_model.Zone
	var emptyConnection template_model.Connection
	zones := linq.SelectSlice(this.ManualZones, func(edit ManualZoneEdit) template_model.Zone {
		return edit.Apply(emptyZone)
	})
	connections := linq.SelectSlice(this.ManualConnections, func(edit ManualConnectionEdit) template_model.Connection {
		return edit.Apply(emptyConnection)
	})
	return zones, connections
}

func (this ManualEditSettings) Equals(other ManualEditSettings) bool {
	return manualEditsEqual(this.ManualZones, other.ManualZones) &&
		manualEditsEqual(this.ManualConnections, other.ManualConnections)
}

func (this ManualEditSettings) Clone() ManualEditSettings {
	return ManualEditSettings{
		ManualZones:       linq.SelectSlice(this.ManualZones, ManualZoneEdit.Clone),
		ManualConnections: linq.SelectSlice(this.ManualConnections, ManualConnectionEdit.Clone),
	}
}

func ToManualEditSettingsModel(entity editor_state.ManualEditSettings) ManualEditSettings {
	return ManualEditSettings{
		ManualZones: linq.SelectSlice(entity.Zones, func(edit editor_state.ManualZoneEdit) ManualZoneEdit {
			return ManualZoneEdit{ManualZoneEdit: edit}
		}),
		ManualConnections: linq.SelectSlice(entity.Connections,
			func(edit editor_state.ManualConnectionEdit) ManualConnectionEdit {
				return ManualConnectionEdit{ManualConnectionEdit: edit}
			}),
	}
}

func ToManualEditSettingsEntity(model ManualEditSettings) editor_state.ManualEditSettings {
	return editor_state.ManualEditSettings{
		Zones: linq.SelectSlice(model.ManualZones, func(edit ManualZoneEdit) editor_state.ManualZoneEdit {
			return edit.ManualZoneEdit
		}),
		Connections: linq.SelectSlice(model.ManualConnections,
			func(edit ManualConnectionEdit) editor_state.ManualConnectionEdit {
				return edit.ManualConnectionEdit
			}),
	}
}

func manualEditsEqual[Edit any](left, right []Edit) bool {
	if len(left) == 0 || len(right) == 0 {
		return len(left) == len(right)
	}

	return reflect.DeepEqual(left, right)
}
