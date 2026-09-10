package connection_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model/template_variant_model"
	"github.com/stretchr/testify/assert"
)

func TestWhenConnectionTypeIsChecked_PortalIsMatchedCaseInsensitively(t *testing.T) {
	t.Parallel()
	for caseName, testCase := range map[string]struct {
		connectionType string
		expected       bool
	}{
		"WhenTypeIsPortal_ReturnsTrue":          {connectionType: "Portal", expected: true},
		"WhenTypeIsLowerCasePortal_ReturnsTrue": {connectionType: "portal", expected: true},
		"WhenTypeIsUpperCasePortal_ReturnsTrue": {connectionType: "PORTAL", expected: true},
		"WhenTypeIsMixedCasePortal_ReturnsTrue": {connectionType: "PoRtAl", expected: true},
		"WhenTypeIsDirect_ReturnsFalse":         {connectionType: "Direct", expected: false},
		"WhenTypeIsDefault_ReturnsFalse":        {connectionType: "Default", expected: false},
		"WhenTypeIsProximity_ReturnsFalse":      {connectionType: "Proximity", expected: false},
		"WhenTypeIsGladiatorArena_ReturnsFalse": {connectionType: "GladiatorArena", expected: false},
		"WhenTypeIsEmpty_ReturnsFalse":          {connectionType: "", expected: false},
		"WhenTypeIsPaddedPortal_ReturnsFalse":   {connectionType: " Portal", expected: false},
	} {
		t.Run(caseName, func(t *testing.T) {
			t.Parallel()
			// Arrange
			connection := template_variant_model.Connection{ConnectionType: testCase.connectionType}

			// Act
			result := connection.IsExplicitPortal()

			// Assert
			assert.Equal(t, testCase.expected, result)
		})
	}
}

func TestWhenOnlyPlacementRulesArePortal_ReturnsFalse(t *testing.T) {
	t.Parallel()
	for caseName, connectionType := range map[string]string{
		"WhenTypeIsDirect_ReturnsFalse": "Direct",
		"WhenTypeIsEmpty_ReturnsFalse":  "",
	} {
		t.Run(caseName, func(t *testing.T) {
			t.Parallel()
			// Arrange
			connection := template_variant_model.Connection{
				ConnectionType:           connectionType,
				PortalPlacementRulesFrom: newPortalPlacementRules(),
				PortalPlacementRulesTo:   newPortalPlacementRules(),
			}

			// Act
			result := connection.IsExplicitPortal()

			// Assert
			assert.False(t, result)
		})
	}
}

func TestWhenRoadFlagIsSet_ExplicitPortalDetectionIsUnaffected(t *testing.T) {
	t.Parallel()
	for caseName, road := range map[string]bool{
		"WhenRoadIsTrue_ReturnsTrue":  true,
		"WhenRoadIsFalse_ReturnsTrue": false,
	} {
		t.Run(caseName, func(t *testing.T) {
			t.Parallel()
			// Arrange
			roadFlag := road
			connection := template_variant_model.Connection{
				ConnectionType: "Portal",
				Road:           &roadFlag,
			}

			// Act
			result := connection.IsExplicitPortal()

			// Assert
			assert.True(t, result)
		})
	}
}
