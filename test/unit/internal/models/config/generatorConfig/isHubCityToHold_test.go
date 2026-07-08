package generatorConfig_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/stretchr/testify/assert"
)

func TestWhenTopologyAndCityHoldCombinationsVary_ReportsHubCityToHoldAccordingly(t *testing.T) {
	testCases := []struct {
		subtestName string
		mutate      func(configuration *config.GeneratorConfig)
		expected    bool
	}{
		{
			"WhenHubTopologyHasCityHoldFlag_ReturnsTrue",
			func(configuration *config.GeneratorConfig) {
				configuration.Topology = config.TopologyHubAndSpoke
				configuration.GameEndConditions.CityHold = true
			},
			true,
		},
		{
			"WhenHubTopologyHasCityHoldVictoryCondition_ReturnsTrue",
			func(configuration *config.GeneratorConfig) {
				configuration.Topology = config.TopologyHubAndSpoke
				configuration.GameEndConditions.VictoryCondition = "win_condition_5"
			},
			true,
		},
		{
			"WhenTopologyIsNotHubDespiteCityHold_ReturnsFalse",
			func(configuration *config.GeneratorConfig) {
				configuration.Topology = config.TopologyRing
				configuration.GameEndConditions.CityHold = true
			},
			false,
		},
		{
			"WhenHubTopologyHasNoCityHoldMode_ReturnsFalse",
			func(configuration *config.GeneratorConfig) {
				configuration.Topology = config.TopologyHubAndSpoke
			},
			false,
		},
		{
			"WhenHubTopologyHasNilGameEndConditions_ReturnsFalse",
			func(configuration *config.GeneratorConfig) {
				configuration.Topology = config.TopologyHubAndSpoke
				configuration.GameEndConditions = nil
			},
			false,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.subtestName, func(t *testing.T) {
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
