package bonusPresetType_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/config_helpers"
	"github.com/stretchr/testify/assert"
)

func TestWhenPresetTypeIsNamed_ReturnsItsName(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		subtestName string
		presetType  editor_state.BonusPresetType
		expected    string
	}{
		{"WhenTypeIsTownPortalFree_ReturnsTownPortalFree", editor_state.BonusTownPortalFree, "TownPortalFree"},
		{"WhenTypeIsSpell_ReturnsSpell", editor_state.BonusSpell, "Spell"},
		{"WhenTypeIsUnitMultiplier_ReturnsUnitMultiplier", editor_state.BonusUnitMultiplier, "UnitMultiplier"},
		{"WhenTypeIsMovementBonus_ReturnsMovementBonus", editor_state.BonusMovementBonus, "MovementBonus"},
		{"WhenTypeIsStartingItem_ReturnsStartingItem", editor_state.BonusStartingItem, "StartingItem"},
		{"WhenTypeIsStartingGold_ReturnsStartingGold", editor_state.BonusStartingGold, "StartingGold"},
		{"WhenTypeIsStartingGems_ReturnsStartingGems", editor_state.BonusStartingGems, "StartingGems"},
		{"WhenTypeIsStartingCrystals_ReturnsStartingCrystals", editor_state.BonusStartingCrystals, "StartingCrystals"},
		{"WhenTypeIsStartingMercury_ReturnsStartingMercury", editor_state.BonusStartingMercury, "StartingMercury"},
		{"WhenTypeIsStartingWood_ReturnsStartingWood", editor_state.BonusStartingWood, "StartingWood"},
		{"WhenTypeIsStartingOre_ReturnsStartingOre", editor_state.BonusStartingOre, "StartingOre"},
	}
	for _, testCase := range testCases {
		t.Run(testCase.subtestName, func(t *testing.T) {
			t.Parallel()
			// Arrange

			// Act
			actual := config_helpers.GetString(testCase.presetType)

			// Assert
			assert.Equal(t, testCase.expected, actual)
		})
	}
}

func TestWhenPresetTypeIsUnknown_ReturnsNumericString(t *testing.T) {
	t.Parallel()
	// Arrange
	unknown := editor_state.BonusPresetType(99)

	// Act
	actual := config_helpers.GetString(unknown)

	// Assert
	assert.Equal(t, "99", actual)
}
