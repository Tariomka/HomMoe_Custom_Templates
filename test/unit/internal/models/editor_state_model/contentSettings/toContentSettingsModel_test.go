package contentSettings_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/stretchr/testify/assert"
)

func TestWhenTheGroupIsPersisted_TheTypedListsAreCarried(t *testing.T) {
	t.Parallel()
	// Arrange
	entity := allFieldsContentSettings()

	// Act
	model := editor_state_model.ToContentSettingsModel(entity)

	// Assert
	assert.Equal(
		t,
		[]any{entity.BannedItems, entity.BannedMagics, entity.ValueOverrides},
		[]any{
			editor_state_model.ToBanEntryEntities(model.BannedItems),
			editor_state_model.ToBanEntryEntities(model.BannedMagics),
			editor_state_model.ToValueOverrideEntryEntities(model.ValueOverrides),
		})
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

func TestWhenAnEmptyGroupIsPersisted_PlayerRowsComeFromTheDefaultsAndTheRestStaysNil(t *testing.T) {
	t.Parallel()
	// Arrange
	entity := editor_state.ContentSettings{ContentDefaultsVersion: 1}

	// Act
	model := editor_state_model.ToContentSettingsModel(entity)

	// Assert
	assert.Equal(t,
		editor_state_model.ContentSettings{PlayerZoneContentRows: editor_state_model.GetDefaultPlayerZoneContentRows()},
		model)
}

func TestWhenPlayerRowsWereCleared_TheyStayCleared(t *testing.T) {
	t.Parallel()
	// Arrange
	entity := editor_state.ContentSettings{
		ContentDefaultsVersion: 1,
		PlayerZoneContentRows:  []editor_state.ZoneContentRow{},
	}

	// Act
	model := editor_state_model.ToContentSettingsModel(entity)

	// Assert
	assert.Empty(t, model.PlayerZoneContentRows)
}

func TestWhenTheDefaultsVersionIsUnknown_OmittedPlayerRowsUseTheCurrentDefaults(t *testing.T) {
	t.Parallel()
	// Arrange
	entity := editor_state.ContentSettings{ContentDefaultsVersion: 0}

	// Act
	model := editor_state_model.ToContentSettingsModel(entity)

	// Assert
	assert.Equal(t, editor_state_model.GetDefaultPlayerZoneContentRows(), model.PlayerZoneContentRows)
}

// allFieldsContentSettings populates every field of the group so a member the
// converter forgets to carry cannot pass unnoticed.
func allFieldsContentSettings() editor_state.ContentSettings {
	return editor_state.ContentSettings{
		ContentDefaultsVersion: 1,
		BannedItems:            []editor_state.BanEntry{{Sid: "banned_item"}},
		BannedMagics:           []editor_state.BanEntry{{Sid: "banned_magic"}},
		ValueOverrides:         []editor_state.ValueOverrideEntry{{Sid: "sid", GuardValue: 1234}},
		Bonuses:                []editor_state.BonusEntry{{PresetType: editor_state.BonusStartingGold, Param: "500"}},

		PlayerZoneContentRows:    []editor_state.ZoneContentRow{{Sid: "player_row", Count: 1}},
		LowestNeutralContentRows: []editor_state.ZoneContentRow{{Sid: "lowest_row", Count: 2}},
		LowNeutralContentRows:    []editor_state.ZoneContentRow{{Sid: "low_row", Count: 3}},
		MediumNeutralContentRows: []editor_state.ZoneContentRow{{Sid: "medium_row", Count: 4}},
		HighNeutralContentRows:   []editor_state.ZoneContentRow{{Sid: "high_row", Count: 5}},
		HubZoneContentRows:       []editor_state.ZoneContentRow{{Sid: "hub_row", Count: 6}},
	}
}
