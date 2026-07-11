package zoneEditor_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/connection_editor"
	"github.com/stretchr/testify/assert"
)

func TestWhenZoneIsChecked_ProtectsOnlyPlayerZones(t *testing.T) {
	t.Parallel()
	playerZoneNames := map[string]bool{"Spawn-A": true, "Spawn-B": true}

	testCases := []struct {
		name     string
		zoneName string
		expected bool
	}{
		{"WhenZoneIsPlayerSpawn_ReturnsFalse", "Spawn-A", false},
		{"WhenZoneIsNeutral_ReturnsTrue", "Neutral-C", true},
		{"WhenZoneIsHub_ReturnsTrue", "Hub", true},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			// Arrange

			// Act
			canDelete := connection_editor.CanDeleteZone(testCase.zoneName, playerZoneNames)

			// Assert
			assert.Equal(t, testCase.expected, canDelete)
		})
	}
}
