package manualReapply_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/connection_editor"
	"github.com/stretchr/testify/assert"
)

func TestWhenZoneNameIsChecked_RecognizesOnlyNeutralPrefix(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name     string
		zoneName string
		expected bool
	}{
		{"WhenNameHasNeutralPrefix_ReturnsTrue", "Neutral-G", true},
		{"WhenNameIsSpawnZone_ReturnsFalse", "Spawn-A", false},
		{"WhenNameIsHub_ReturnsFalse", "Hub", false},
		{"WhenNameIsEmpty_ReturnsFalse", "", false},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			// Arrange

			// Act
			isNeutral := connection_editor.IsNeutralZoneName(testCase.zoneName)

			// Assert
			assert.Equal(t, testCase.expected, isNeutral)
		})
	}
}
