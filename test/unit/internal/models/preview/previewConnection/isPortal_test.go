package previewConnection_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/preview"
	"github.com/stretchr/testify/assert"
)

func TestWhenConnectionTypeIsPortal_ReturnsTrue(t *testing.T) {
	t.Parallel()
	// Arrange
	connection := preview.Connection{Type: preview.ConnectionTypePortal}

	// Act
	isPortal := connection.IsPortal()

	// Assert
	assert.True(t, isPortal)
}

func TestWhenConnectionTypeIsNotPortal_ReturnsFalse(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		subtestName    string
		connectionType preview.ConnectionType
	}{
		{"WhenConnectionTypeIsDirect_ReturnsFalse", preview.ConnectionTypeDirect},
		{"WhenConnectionTypeIsGladiatorArena_ReturnsFalse", preview.ConnectionTypeGladiatorArena},
		{"WhenConnectionTypeIsProximity_ReturnsFalse", preview.ConnectionTypeProximity},
	}
	for _, testCase := range testCases {
		t.Run(testCase.subtestName, func(t *testing.T) {
			t.Parallel()
			// Arrange
			connection := preview.Connection{Type: testCase.connectionType}

			// Act
			isPortal := connection.IsPortal()

			// Assert
			assert.False(t, isPortal)
		})
	}
}
