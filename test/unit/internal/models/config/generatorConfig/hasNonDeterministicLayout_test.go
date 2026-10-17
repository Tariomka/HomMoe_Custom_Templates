package generatorConfig_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/stretchr/testify/assert"
)

func TestWhenTopologyAndRandomPortalsVary_ReportsNonDeterministicLayoutAccordingly(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		subtestName   string
		topology      config.MapTopology
		randomPortals bool
		expected      bool
	}{
		{"WhenTopologyIsRandom_ReturnsTrue", config.TopologyRandom, false, true},
		{"WhenTopologyIsRandomWithRandomPortals_ReturnsTrue", config.TopologyRandom, true, true},
		{"WhenSquareTopologyHasRandomPortals_ReturnsTrue", config.TopologySquare, true, true},
		{"WhenSquareTopologyHasNoRandomPortals_ReturnsFalse", config.TopologySquare, false, false},
		{"WhenCirclesTopologyHasNoRandomPortals_ReturnsFalse", config.TopologyCircles, false, false},
	}
	for _, testCase := range testCases {
		t.Run(testCase.subtestName, func(t *testing.T) {
			t.Parallel()
			// Arrange
			configuration := config.NewGeneratorConfig()
			configuration.Topology = testCase.topology
			configuration.RandomPortals = testCase.randomPortals

			// Act
			nonDeterministic := configuration.HasNonDeterministicLayout()

			// Assert
			assert.Equal(t, testCase.expected, nonDeterministic)
		})
	}
}
