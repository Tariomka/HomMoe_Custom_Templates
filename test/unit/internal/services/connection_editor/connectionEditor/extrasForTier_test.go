package connectionEditor_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/connection_editor"
	"github.com/stretchr/testify/assert"
)

func TestWhenTierIsGiven_ReturnsGeneratorDefaultExtra(t *testing.T) {
	testCases := []struct {
		name     string
		tier     connection_editor.ZoneTier
		expected int
	}{
		{"WhenTierIsBronze_DefaultsTo15000", connection_editor.ZoneTierBronze, 15000},
		{"WhenTierIsSilver_DefaultsTo20000", connection_editor.ZoneTierSilver, 20000},
		{"WhenTierIsGold_DefaultsTo25000", connection_editor.ZoneTierGold, 25000},
		{"WhenTierIsPlayerToPlayer_DefaultsTo30000", connection_editor.ZoneTierPlayerToPlayer, 30000},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// Arrange
			expected := []connection_editor.GuardPresetExtra{
				{Label: "Generator Default", Value: testCase.expected},
			}

			// Act
			extras := connection_editor.ExtrasForTier(testCase.tier)

			// Assert
			assert.Equal(t, expected, extras)
		})
	}
}
