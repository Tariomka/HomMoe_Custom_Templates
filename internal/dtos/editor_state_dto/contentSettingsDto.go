package editor_state_dto

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
)

type ContentSettingsDto struct {
	BannedItems        string
	BannedMagics       string
	ValueOverridesText string
	Bonuses            []config.BonusEntry

	PlayerZoneContentRows    []models.ZoneContentRow
	LowestNeutralContentRows []models.ZoneContentRow
	LowNeutralContentRows    []models.ZoneContentRow
	MediumNeutralContentRows []models.ZoneContentRow
	HighNeutralContentRows   []models.ZoneContentRow
	HubZoneContentRows       []models.ZoneContentRow
}
