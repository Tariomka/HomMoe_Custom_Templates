package bonusPresetType_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/config/config_inner"
	"github.com/stretchr/testify/assert"
)

func TestWhenPresetTypeIsNamed_ReturnsItsName(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		subtestName string
		presetType  config_inner.BonusPresetType
		expected    string
	}{
		{"WhenTypeIsTownPortalFree_ReturnsTownPortalFree", config_inner.BonusTownPortalFree, "TownPortalFree"},
		{"WhenTypeIsSpell_ReturnsSpell", config_inner.BonusSpell, "Spell"},
		{"WhenTypeIsUnitMultiplier_ReturnsUnitMultiplier", config_inner.BonusUnitMultiplier, "UnitMultiplier"},
		{"WhenTypeIsMovementBonus_ReturnsMovementBonus", config_inner.BonusMovementBonus, "MovementBonus"},
		{"WhenTypeIsStartingItem_ReturnsStartingItem", config_inner.BonusStartingItem, "StartingItem"},
		{"WhenTypeIsStartingGold_ReturnsStartingGold", config_inner.BonusStartingGold, "StartingGold"},
		{"WhenTypeIsStartingGems_ReturnsStartingGems", config_inner.BonusStartingGems, "StartingGems"},
		{"WhenTypeIsStartingCrystals_ReturnsStartingCrystals", config_inner.BonusStartingCrystals, "StartingCrystals"},
		{"WhenTypeIsStartingMercury_ReturnsStartingMercury", config_inner.BonusStartingMercury, "StartingMercury"},
		{"WhenTypeIsStartingWood_ReturnsStartingWood", config_inner.BonusStartingWood, "StartingWood"},
		{"WhenTypeIsStartingOre_ReturnsStartingOre", config_inner.BonusStartingOre, "StartingOre"},
	}
	for _, testCase := range testCases {
		t.Run(testCase.subtestName, func(t *testing.T) {
			t.Parallel()
			// Arrange

			// Act
			actual := testCase.presetType.String()

			// Assert
			assert.Equal(t, testCase.expected, actual)
		})
	}
}

func TestWhenPresetTypeIsUnknown_ReturnsNumericString(t *testing.T) {
	t.Parallel()
	// Arrange
	unknown := config_inner.BonusPresetType(99)

	// Act
	actual := unknown.String()

	// Assert
	assert.Equal(t, "99", actual)
}
