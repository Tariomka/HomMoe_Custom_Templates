package bonusPresetType_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/config_helpers"
	"github.com/stretchr/testify/assert"
)

func TestWhenPresetTypeVaries_ReportsResourceKindsOnly(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		subtestName string
		presetType  editor_state.BonusPresetType
		expected    bool
	}{
		{"WhenTypeIsTownPortalFree_ReturnsFalse", editor_state.BonusTownPortalFree, false},
		{"WhenTypeIsSpell_ReturnsFalse", editor_state.BonusSpell, false},
		{"WhenTypeIsUnitMultiplier_ReturnsFalse", editor_state.BonusUnitMultiplier, false},
		{"WhenTypeIsMovementBonus_ReturnsFalse", editor_state.BonusMovementBonus, false},
		{"WhenTypeIsStartingItem_ReturnsFalse", editor_state.BonusStartingItem, false},
		{"WhenTypeIsStartingGold_ReturnsTrue", editor_state.BonusStartingGold, true},
		{"WhenTypeIsStartingGems_ReturnsTrue", editor_state.BonusStartingGems, true},
		{"WhenTypeIsStartingCrystals_ReturnsTrue", editor_state.BonusStartingCrystals, true},
		{"WhenTypeIsStartingMercury_ReturnsTrue", editor_state.BonusStartingMercury, true},
		{"WhenTypeIsStartingWood_ReturnsTrue", editor_state.BonusStartingWood, true},
		{"WhenTypeIsStartingOre_ReturnsTrue", editor_state.BonusStartingOre, true},
	}
	for _, testCase := range testCases {
		t.Run(testCase.subtestName, func(t *testing.T) {
			t.Parallel()
			// Arrange

			// Act
			actual := config_helpers.IsResource(testCase.presetType)

			// Assert
			assert.Equal(t, testCase.expected, actual)
		})
	}
}
