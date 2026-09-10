package connection_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model/template_variant_model"
	"github.com/stretchr/testify/assert"
)

func TestWhenRoadFlagIsExplicit_ItWinsRegardlessOfType(t *testing.T) {
	t.Parallel()
	for caseName, testCase := range map[string]struct {
		connectionType string
		road           bool
		expected       bool
	}{
		"WhenRoadIsFalseOnDirect_ReturnsFalse":      {connectionType: "Direct", road: false, expected: false},
		"WhenRoadIsFalseOnPortal_ReturnsFalse":      {connectionType: "Portal", road: false, expected: false},
		"WhenRoadIsFalseOnLowerPortal_ReturnsFalse": {connectionType: "portal", road: false, expected: false},
		"WhenRoadIsTrueOnDirect_ReturnsTrue":        {connectionType: "Direct", road: true, expected: true},
		"WhenRoadIsTrueOnPortal_ReturnsTrue":        {connectionType: "Portal", road: true, expected: true},
		"WhenRoadIsTrueOnEmptyType_ReturnsTrue":     {connectionType: "", road: true, expected: true},
	} {
		t.Run(caseName, func(t *testing.T) {
			t.Parallel()
			// Arrange
			road := testCase.road
			connection := template_variant_model.Connection{
				ConnectionType: testCase.connectionType,
				Road:           &road,
			}

			// Act
			result := connection.HasRoad()

			// Assert
			assert.Equal(t, testCase.expected, result)
		})
	}
}

func TestWhenRoadFlagIsNil_OnlyExplicitPortalIsRoaded(t *testing.T) {
	t.Parallel()
	for caseName, testCase := range map[string]struct {
		connectionType string
		expected       bool
	}{
		"WhenTypeIsPortal_ReturnsTrue":          {connectionType: "Portal", expected: true},
		"WhenTypeIsLowerCasePortal_ReturnsTrue": {connectionType: "portal", expected: true},
		"WhenTypeIsMixedCasePortal_ReturnsTrue": {connectionType: "PoRtAl", expected: true},
		"WhenTypeIsDirect_ReturnsFalse":         {connectionType: "Direct", expected: false},
		"WhenTypeIsDefault_ReturnsFalse":        {connectionType: "Default", expected: false},
		"WhenTypeIsProximity_ReturnsFalse":      {connectionType: "Proximity", expected: false},
		"WhenTypeIsGladiatorArena_ReturnsFalse": {connectionType: "GladiatorArena", expected: false},
		"WhenTypeIsEmpty_ReturnsFalse":          {connectionType: "", expected: false},
	} {
		t.Run(caseName, func(t *testing.T) {
			t.Parallel()
			// Arrange
			connection := template_variant_model.Connection{ConnectionType: testCase.connectionType}

			// Act
			result := connection.HasRoad()

			// Assert
			assert.Equal(t, testCase.expected, result)
		})
	}
}

func TestWhenRoadFlagIsNilAndOnlyPlacementRulesArePortal_ReturnsFalse(t *testing.T) {
	t.Parallel()
	// Arrange
	connection := template_variant_model.Connection{
		ConnectionType:           "Direct",
		PortalPlacementRulesFrom: newPortalPlacementRules(),
		PortalPlacementRulesTo:   newPortalPlacementRules(),
	}

	// Act
	result := connection.HasRoad()

	// Assert
	assert.False(t, result)
}
