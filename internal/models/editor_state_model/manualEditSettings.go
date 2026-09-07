package editor_state_model

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
)

type ManualEditSettings struct {
	ManualZones       []template_model.Zone
	ManualConnections []ManualConnectionSave
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
