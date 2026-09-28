package connection_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model/template_variant_model"
	"github.com/stretchr/testify/assert"
)

func TestWhenOnlyTheTypeIsConsidered_PortalIsMatchedCaseInsensitively(t *testing.T) {
	t.Parallel()
	for caseName, testCase := range map[string]struct {
		connectionType string
		expected       bool
	}{
		"WhenTypeIsPortal_ReturnsTrue":          {connectionType: "Portal", expected: true},
		"WhenTypeIsLowerCasePortal_ReturnsTrue": {connectionType: "portal", expected: true},
		"WhenTypeIsDirect_ReturnsFalse":         {connectionType: "Direct", expected: false},
		"WhenTypeIsProximity_ReturnsFalse":      {connectionType: "Proximity", expected: false},
		"WhenTypeIsEmpty_ReturnsFalse":          {connectionType: "", expected: false},
	} {
		t.Run(caseName, func(t *testing.T) {
			t.Parallel()
			// Arrange
			connection := template_variant_model.Connection{ConnectionType: testCase.connectionType}

			// Act
			result := connection.IsEffectivePortal()

			// Assert
			assert.Equal(t, testCase.expected, result)
		})
	}
}

func TestWhenADirectConnectionCarriesSourcePlacementRules_ReturnsTrue(t *testing.T) {
	t.Parallel()
	// Arrange
	connection := template_variant_model.Connection{
		ConnectionType:           "Direct",
		PortalPlacementRulesFrom: newPortalPlacementRules(),
	}

	// Act
	result := connection.IsEffectivePortal()

	// Assert
	assert.True(t, result)
}

func TestWhenADirectConnectionCarriesTargetPlacementRules_ReturnsTrue(t *testing.T) {
	t.Parallel()
	// Arrange
	connection := template_variant_model.Connection{
		ConnectionType:         "Direct",
		PortalPlacementRulesTo: newPortalPlacementRules(),
	}

	// Act
	result := connection.IsEffectivePortal()

	// Assert
	assert.True(t, result)
}

func TestWhenAnUntypedConnectionCarriesPlacementRules_ReturnsTrue(t *testing.T) {
	t.Parallel()
	// Arrange
	connection := template_variant_model.Connection{PortalPlacementRulesFrom: newPortalPlacementRules()}

	// Act
	result := connection.IsEffectivePortal()

	// Assert
	assert.True(t, result)
}
