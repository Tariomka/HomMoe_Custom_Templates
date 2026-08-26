package bonusEntry_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/config_helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenEntriesAreIdentical_ProducesSameHash(t *testing.T) {
	t.Parallel()
	// Arrange
	entry := editor_state_model.BonusEntry{
		PresetType:     editor_state_model.BonusSpell,
		ReceiverFilter: "start_hero",
		Param:          gofakeit.Word(),
		Param2:         "1",
	}
	duplicate := entry

	// Act
	actual := config_helpers.GetHash(entry)

	// Assert
	assert.Equal(t, config_helpers.GetHash(duplicate), actual)
}

func TestWhenEntriesDifferInParam_ProducesDifferentHashes(t *testing.T) {
	t.Parallel()
	// Arrange
	entry := editor_state_model.BonusEntry{
		PresetType:     editor_state_model.BonusStartingGold,
		ReceiverFilter: "all_heroes",
		Param:          "500",
	}
	other := entry
	other.Param = "700"

	// Act
	actual := config_helpers.GetHash(entry)

	// Assert
	assert.NotEqual(t, config_helpers.GetHash(other), actual)
}

func TestWhenEntriesDifferInPresetType_ProducesDifferentHashes(t *testing.T) {
	t.Parallel()
	// Arrange
	entry := editor_state_model.BonusEntry{
		PresetType:     editor_state_model.BonusStartingWood,
		ReceiverFilter: "start_hero",
		Param:          "7",
	}
	other := entry
	other.PresetType = editor_state_model.BonusStartingOre

	// Act
	actual := config_helpers.GetHash(entry)

	// Assert
	assert.NotEqual(t, config_helpers.GetHash(other), actual)
}

func TestWhenHashIsComputed_ReturnsSha256DigestLength(t *testing.T) {
	t.Parallel()
	// Arrange
	entry := editor_state_model.BonusEntry{
		PresetType:     editor_state_model.BonusTownPortalFree,
		ReceiverFilter: gofakeit.Word(),
	}

	// Act
	actual := config_helpers.GetHash(entry)

	// Assert
	assert.Len(t, actual, 32)
}
