package bonusPresetType_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/config/config_inner"
	"github.com/stretchr/testify/assert"
)

func TestWhenPresetTypeVaries_ReportsResourceKindsOnly(t *testing.T) {
	testCases := []struct {
		subtestName string
		presetType  config_inner.BonusPresetType
		expected    bool
	}{
		{"WhenTypeIsTownPortalFree_ReturnsFalse", config_inner.BonusTownPortalFree, false},
		{"WhenTypeIsSpell_ReturnsFalse", config_inner.BonusSpell, false},
		{"WhenTypeIsUnitMultiplier_ReturnsFalse", config_inner.BonusUnitMultiplier, false},
		{"WhenTypeIsMovementBonus_ReturnsFalse", config_inner.BonusMovementBonus, false},
		{"WhenTypeIsStartingItem_ReturnsFalse", config_inner.BonusStartingItem, false},
		{"WhenTypeIsStartingGold_ReturnsTrue", config_inner.BonusStartingGold, true},
		{"WhenTypeIsStartingGems_ReturnsTrue", config_inner.BonusStartingGems, true},
		{"WhenTypeIsStartingCrystals_ReturnsTrue", config_inner.BonusStartingCrystals, true},
		{"WhenTypeIsStartingMercury_ReturnsTrue", config_inner.BonusStartingMercury, true},
		{"WhenTypeIsStartingWood_ReturnsTrue", config_inner.BonusStartingWood, true},
		{"WhenTypeIsStartingOre_ReturnsTrue", config_inner.BonusStartingOre, true},
	}
	for _, testCase := range testCases {
		t.Run(testCase.subtestName, func(t *testing.T) {
			// Arrange - preset type provided by the test case.

			// Act
			actual := testCase.presetType.IsResource()

			// Assert
			assert.Equal(t, testCase.expected, actual)
		})
	}
}
