package editor_state_model

import (
	"reflect"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
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
		ManualZones:       ToManualZoneModels(entity.ManualZones),
		ManualConnections: ToManualConnectionSaveModels(entity.ManualConnections),
	}
}

func ToManualEditSettingsEntity(model ManualEditSettings) editor_state.ManualEditSettings {
	return editor_state.ManualEditSettings{
		ManualZones:       ToManualZoneSaveEntities(model.ManualZones),
		ManualConnections: ToManualConnectionSaveEntities(model.ManualConnections),
	}
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
