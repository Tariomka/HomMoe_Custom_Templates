package generatorConfig_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/stretchr/testify/assert"
)

func TestWhenTopologyAndCityHoldCombinationsVary_ReportsHubCityToHoldAccordingly(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		subtestName string
		mutate      func(configuration *config.GeneratorConfig)
		expected    bool
	}{
		{
			"WhenHubTopologyHasCityHoldFlag_ReturnsTrue",
			func(configuration *config.GeneratorConfig) {
				configuration.Topology = config.TopologyGeometricHub
				configuration.GameEndConditions.CityHold = true
			},
			true,
		},
		{
			"WhenHubTopologyHasCityHoldVictoryCondition_ReturnsTrue",
			func(configuration *config.GeneratorConfig) {
				configuration.Topology = config.TopologyGeometricHub
				configuration.GameEndConditions.VictoryCondition = "win_condition_5"
			},
			true,
		},
		{
			"WhenTopologyIsNotHubDespiteCityHold_ReturnsFalse",
			func(configuration *config.GeneratorConfig) {
				configuration.Topology = config.TopologyCircles
				configuration.GameEndConditions.CityHold = true
			},
			false,
		},
		{
			"WhenTopologyIsTheRetiredHubDespiteCityHold_ReturnsFalse",
			func(configuration *config.GeneratorConfig) {
				configuration.Topology = config.MapTopology("HubAndSpoke")
				configuration.GameEndConditions.CityHold = true
			},
			false,
		},
		{
			"WhenHubTopologyHasNoCityHoldMode_ReturnsFalse",
			func(configuration *config.GeneratorConfig) {
				configuration.Topology = config.TopologyGeometricHub
			},
			false,
		},
		{
			"WhenHubTopologyHasNilGameEndConditions_ReturnsFalse",
			func(configuration *config.GeneratorConfig) {
				configuration.Topology = config.TopologyGeometricHub
				configuration.GameEndConditions = nil
			},
			false,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.subtestName, func(t *testing.T) {
			t.Parallel()
			// Arrange
			configuration := config.NewGeneratorConfig()
			testCase.mutate(configuration)

			// Act
			actual := configuration.IsHubCityToHold()

			// Assert
			assert.Equal(t, testCase.expected, actual)
		})
	}
}
