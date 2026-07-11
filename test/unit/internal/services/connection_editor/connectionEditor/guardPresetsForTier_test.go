package connectionEditor_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/connection_editor"
	"github.com/stretchr/testify/assert"
)

func TestWhenTierIsGiven_ReturnsMatchingGuardPresetRow(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name     string
		tier     connection_editor.ZoneTier
		expected [5]int
	}{
		{"WhenTierIsBronze_ReturnsBronzeRow", connection_editor.ZoneTierBronze, [5]int{3000, 6000, 9000, 12000, 16000}},
		{
			"WhenTierIsSilver_ReturnsSilverRow",
			connection_editor.ZoneTierSilver,
			[5]int{18000, 21000, 24000, 27000, 30000},
		},
		{"WhenTierIsGold_ReturnsGoldRow", connection_editor.ZoneTierGold, [5]int{36000, 42000, 48000, 54000, 60000}},
		{
			"WhenTierIsPlayerToPlayer_ReturnsPlayerToPlayerRow",
			connection_editor.ZoneTierPlayerToPlayer,
			[5]int{10000, 22000, 34000, 46000, 58000},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			// Arrange

			// Act
			presets := connection_editor.GuardPresetsForTier(testCase.tier)

			// Assert
			assert.Equal(t, testCase.expected, presets)
		})
	}
}
