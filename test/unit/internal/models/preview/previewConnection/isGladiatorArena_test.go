package previewConnection_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/preview"
	"github.com/stretchr/testify/assert"
)

func TestWhenConnectionTypeIsGladiatorArena_ReturnsTrue(t *testing.T) {
	t.Parallel()
	// Arrange
	connection := preview.Connection{Type: preview.ConnectionTypeGladiatorArena}

	// Act
	isArena := connection.IsGladiatorArena()

	// Assert
	assert.True(t, isArena)
}

func TestWhenConnectionTypeIsNotGladiatorArena_ReturnsFalse(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		subtestName    string
		connectionType preview.ConnectionType
	}{
		{"WhenConnectionTypeIsDirect_ReturnsFalse", preview.ConnectionTypeDirect},
		{"WhenConnectionTypeIsPortal_ReturnsFalse", preview.ConnectionTypePortal},
		{"WhenConnectionTypeIsProximity_ReturnsFalse", preview.ConnectionTypeProximity},
	}
	for _, testCase := range testCases {
		t.Run(testCase.subtestName, func(t *testing.T) {
			t.Parallel()
			// Arrange
			connection := preview.Connection{Type: testCase.connectionType}

			// Act
			isArena := connection.IsGladiatorArena()

			// Assert
			assert.False(t, isArena)
		})
	}
}
