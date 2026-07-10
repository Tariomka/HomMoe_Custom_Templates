package bonusEntry_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/config/config_inner"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenEntriesAreIdentical_ProducesSameHash(t *testing.T) {
	t.Parallel()
	// Arrange
	entry := config_inner.BonusEntry{
		PresetType:     config_inner.BonusSpell,
		ReceiverFilter: "start_hero",
		Param:          gofakeit.Word(),
		Param2:         "1",
	}
	duplicate := entry

	// Act
	actual := entry.GetHash()

	// Assert
	assert.Equal(t, duplicate.GetHash(), actual)
}

func TestWhenEntriesDifferInParam_ProducesDifferentHashes(t *testing.T) {
	t.Parallel()
	// Arrange
	entry := config_inner.BonusEntry{
		PresetType:     config_inner.BonusStartingGold,
		ReceiverFilter: "all_heroes",
		Param:          "500",
	}
	other := entry
	other.Param = "700"

	// Act
	actual := entry.GetHash()

	// Assert
	assert.NotEqual(t, other.GetHash(), actual)
}

func TestWhenEntriesDifferInPresetType_ProducesDifferentHashes(t *testing.T) {
	t.Parallel()
	// Arrange
	entry := config_inner.BonusEntry{
		PresetType:     config_inner.BonusStartingWood,
		ReceiverFilter: "start_hero",
		Param:          "7",
	}
	other := entry
	other.PresetType = config_inner.BonusStartingOre

	// Act
	actual := entry.GetHash()

	// Assert
	assert.NotEqual(t, other.GetHash(), actual)
}

func TestWhenHashIsComputed_ReturnsSha256DigestLength(t *testing.T) {
	t.Parallel()
	// Arrange
	entry := config_inner.BonusEntry{
		PresetType:     config_inner.BonusTownPortalFree,
		ReceiverFilter: gofakeit.Word(),
	}

	// Act
	actual := entry.GetHash()

	// Assert
	assert.Len(t, actual, 32)
}
