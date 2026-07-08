package generatorConfig_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/stretchr/testify/assert"
)

func TestWhenSeparationConstraintsVary_ReportsHonorabilityAccordingly(t *testing.T) {
	testCases := []struct {
		subtestName string
		mutate      func(configuration *config.GeneratorConfig)
		expected    bool
	}{
		{
			"WhenMinimumSeparationIsZero_ReturnsTrue",
			func(configuration *config.GeneratorConfig) {
				configuration.MinNeutralZonesBetweenPlayers = 0
			},
			true,
		},
		{
			"WhenRandomPortalsAreEnabled_ReturnsFalse",
			func(configuration *config.GeneratorConfig) {
				configuration.MinNeutralZonesBetweenPlayers = 1
				configuration.RandomPortals = true
				configuration.ZoneConfiguration.NeutralZoneCount = 8
			},
			false,
		},
		{
			"WhenRingTopologyHasEnoughNeutrals_ReturnsTrue",
			func(configuration *config.GeneratorConfig) {
				configuration.Topology = config.TopologyRing
				configuration.PlayerCount = 2
				configuration.MinNeutralZonesBetweenPlayers = 2
				configuration.ZoneConfiguration.NeutralZoneCount = 4
			},
			true,
		},
		{
			"WhenRingTopologyLacksNeutrals_ReturnsFalse",
			func(configuration *config.GeneratorConfig) {
				configuration.Topology = config.TopologyRing
				configuration.PlayerCount = 2
				configuration.MinNeutralZonesBetweenPlayers = 2
				configuration.ZoneConfiguration.NeutralZoneCount = 3
			},
			false,
		},
		{
			"WhenCirclesTopologyHasEnoughNeutrals_ReturnsTrue",
			func(configuration *config.GeneratorConfig) {
				configuration.Topology = config.TopologyCircles
				configuration.PlayerCount = 3
				configuration.MinNeutralZonesBetweenPlayers = 1
				configuration.ZoneConfiguration.NeutralZoneCount = 3
			},
			true,
		},
		{
			"WhenChainTopologyHasEnoughNeutralsForGapsBetweenPlayers_ReturnsTrue",
			func(configuration *config.GeneratorConfig) {
				configuration.Topology = config.TopologyChain
				configuration.PlayerCount = 3
				configuration.MinNeutralZonesBetweenPlayers = 2
				configuration.ZoneConfiguration.NeutralZoneCount = 4
			},
			true,
		},
		{
			"WhenChainTopologyLacksNeutrals_ReturnsFalse",
			func(configuration *config.GeneratorConfig) {
				configuration.Topology = config.TopologyChain
				configuration.PlayerCount = 3
				configuration.MinNeutralZonesBetweenPlayers = 2
				configuration.ZoneConfiguration.NeutralZoneCount = 3
			},
			false,
		},
		{
			"WhenHubTopologyRequiresSingleSeparation_ReturnsTrue",
			func(configuration *config.GeneratorConfig) {
				configuration.Topology = config.TopologyHubAndSpoke
				configuration.MinNeutralZonesBetweenPlayers = 1
			},
			true,
		},
		{
			"WhenHubTopologyRequiresDoubleSeparation_ReturnsFalse",
			func(configuration *config.GeneratorConfig) {
				configuration.Topology = config.TopologyHubAndSpoke
				configuration.MinNeutralZonesBetweenPlayers = 2
				configuration.ZoneConfiguration.NeutralZoneCount = 8
			},
			false,
		},
		{
			"WhenSharedWebTopologyRequiresSingleSeparation_ReturnsTrue",
			func(configuration *config.GeneratorConfig) {
				configuration.Topology = config.TopologySharedWeb
				configuration.MinNeutralZonesBetweenPlayers = 1
				configuration.ZoneConfiguration.NeutralZoneCount = 0
			},
			true,
		},
		{
			"WhenSharedWebTopologyRequiresDoubleSeparation_ReturnsFalse",
			func(configuration *config.GeneratorConfig) {
				configuration.Topology = config.TopologySharedWeb
				configuration.MinNeutralZonesBetweenPlayers = 2
				configuration.ZoneConfiguration.NeutralZoneCount = 8
			},
			false,
		},
		{
			"WhenTopologyHasNoSeparationSupport_ReturnsFalse",
			func(configuration *config.GeneratorConfig) {
				configuration.Topology = config.TopologyRandom
				configuration.MinNeutralZonesBetweenPlayers = 1
				configuration.ZoneConfiguration.NeutralZoneCount = 8
			},
			false,
		},
		{
			"WhenAdvancedCountsOverrideSimpleCountWithEnoughNeutrals_ReturnsTrue",
			func(configuration *config.GeneratorConfig) {
				configuration.Topology = config.TopologyRing
				configuration.PlayerCount = 2
				configuration.MinNeutralZonesBetweenPlayers = 2
				configuration.ZoneConfiguration.NeutralZoneCount = 0
				configuration.ZoneConfiguration.Advanced.NeutralLowNoCastleCount = 1
				configuration.ZoneConfiguration.Advanced.NeutralLowCastleCount = 1
				configuration.ZoneConfiguration.Advanced.NeutralMediumNoCastleCount = 1
				configuration.ZoneConfiguration.Advanced.NeutralMediumCastleCount = 1
			},
			true,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.subtestName, func(t *testing.T) {
			// Arrange
			configuration := config.NewGeneratorConfig()
			testCase.mutate(configuration)

			// Act
			actual := configuration.CanHonorNeutralSeparation()

			// Assert
			assert.Equal(t, testCase.expected, actual)
		})
	}
}
