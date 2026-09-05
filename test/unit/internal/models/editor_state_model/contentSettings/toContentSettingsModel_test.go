package contentSettings_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/stretchr/testify/assert"
)

func TestWhenTheGroupIsPersisted_TheScalarsAreCarried(t *testing.T) {
	t.Parallel()
	// Arrange
	entity := allFieldsContentSettings()

	// Act
	model := editor_state_model.ToContentSettingsModel(entity)

	// Assert
	assert.Equal(
		t,
		[]string{entity.BannedItems, entity.BannedMagics, entity.ValueOverridesText},
		[]string{model.BannedItems, model.BannedMagics, model.ValueOverridesText})
}

func TestWhenTheGroupIsPersisted_EveryRowListIsCarried(t *testing.T) {
	t.Parallel()
	// Arrange
	entity := allFieldsContentSettings()

	// Act
	model := editor_state_model.ToContentSettingsModel(entity)

	// Assert
	assert.Equal(
		t,
		[]string{"player_row", "lowest_row", "low_row", "medium_row", "high_row", "hub_row"},
		[]string{
			model.PlayerZoneContentRows[0].Sid,
			model.LowestNeutralContentRows[0].Sid,
			model.LowNeutralContentRows[0].Sid,
			model.MediumNeutralContentRows[0].Sid,
			model.HighNeutralContentRows[0].Sid,
			model.HubZoneContentRows[0].Sid,
		})
}

func TestWhenTheGroupIsPersisted_TheBonusesAreCarried(t *testing.T) {
	t.Parallel()
	// Arrange
	entity := allFieldsContentSettings()

	// Act
	model := editor_state_model.ToContentSettingsModel(entity)

	// Assert
	assert.Equal(t, []editor_state_model.BonusEntry{{BonusEntry: entity.Bonuses[0]}}, model.Bonuses)
}

func TestWhenAnEmptyGroupIsPersisted_TheRowListsStayNil(t *testing.T) {
	t.Parallel()
	// Arrange
	entity := editor_state.ContentSettings{}

	// Act
	model := editor_state_model.ToContentSettingsModel(entity)

	// Assert
	assert.Equal(t, editor_state_model.ContentSettings{}, model)
}

// allFieldsContentSettings populates every field of the group so a member the
// converter forgets to carry cannot pass unnoticed.
func allFieldsContentSettings() editor_state.ContentSettings {
	return editor_state.ContentSettings{
		BannedItems:        "banned_item",
		BannedMagics:       "banned_magic",
		ValueOverridesText: "sid=1234",
		Bonuses:            []editor_state.BonusEntry{{PresetType: editor_state.BonusStartingGold, Param: "500"}},

		PlayerZoneContentRows:    []editor_state.ZoneContentRow{{Sid: "player_row", Count: 1}},
		LowestNeutralContentRows: []editor_state.ZoneContentRow{{Sid: "lowest_row", Count: 2}},
		LowNeutralContentRows:    []editor_state.ZoneContentRow{{Sid: "low_row", Count: 3}},
		MediumNeutralContentRows: []editor_state.ZoneContentRow{{Sid: "medium_row", Count: 4}},
		HighNeutralContentRows:   []editor_state.ZoneContentRow{{Sid: "high_row", Count: 5}},
		HubZoneContentRows:       []editor_state.ZoneContentRow{{Sid: "hub_row", Count: 6}},
	}
}
