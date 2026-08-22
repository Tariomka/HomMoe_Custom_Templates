package editor_state

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
)

// ContentSettings holds the global bans, value overrides, bonuses and the
// mandatory content rows seeded into each zone type.
type ContentSettings struct {
	BannedItems        string              `json:"bannedItems"`
	BannedMagics       string              `json:"bannedMagics"`
	ValueOverridesText string              `json:"valueOverrides"`
	Bonuses            []config.BonusEntry `json:"bonuses"`

	PlayerZoneContentRows    []models.ZoneContentRowSave `json:"playerZoneContentRows,omitempty"`
	LowestNeutralContentRows []models.ZoneContentRowSave `json:"lowestNeutralContentRows,omitempty"`
	LowNeutralContentRows    []models.ZoneContentRowSave `json:"lowNeutralContentRows,omitempty"`
	MediumNeutralContentRows []models.ZoneContentRowSave `json:"mediumNeutralContentRows,omitempty"`
	HighNeutralContentRows   []models.ZoneContentRowSave `json:"highNeutralContentRows,omitempty"`
	HubZoneContentRows       []models.ZoneContentRowSave `json:"hubZoneContentRows,omitempty"`
}
