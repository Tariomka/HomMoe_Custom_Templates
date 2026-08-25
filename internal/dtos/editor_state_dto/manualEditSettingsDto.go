package editor_state_dto

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
)

type ManualEditSettingsDto struct {
	ManualZones       []models.ManualZoneSave
	ManualConnections []models.ManualConnectionSave
}
