package topologyServiceLookup_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
)

func TestWhenTopologyIsMapped_ReturnsCreator(t *testing.T) {
	t.Parallel()
	mappedTopologies := []config.MapTopology{
		config.TopologyHubAndSpoke,
		config.TopologyGeometricHub,
		config.TopologyChain,
		config.TopologySharedWeb,
		config.TopologyRandom,
		config.TopologyCircles,
		config.TopologySquare,
		config.TopologyGeometric,
		config.TopologyCross,
		config.TopologyFractal,
	}

	for _, mapTopology := range mappedTopologies {
		t.Run(string(mapTopology), func(t *testing.T) {
			t.Parallel()
			// Arrange
			lookup := test_helpers.NewTopologyServiceLookup(test_helpers.NewZoneFactories())

			// Act
			creator := lookup.Resolve(mapTopology)

			// Assert
			assert.NotNil(t, creator)
		})
	}
}

func TestWhenTopologyIsHubAndSpoke_ResolvesHubTopology(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyHubAndSpoke
	playerLabels := []string{"A", "B", "C"}
	lookup := test_helpers.NewTopologyServiceLookup(test_helpers.NewZoneFactories())

	// Act
	variant := lookup.Resolve(configuration.Topology)(
		*configuration, playerLabels, neutral_zone.Plans{},
		test_helpers.NewGenerationTuning(configuration, len(playerLabels)), "")

	// Assert
	assert.True(t, hasZoneNamed(variant, "Hub"))
}

func TestWhenTopologyIsUnmapped_ResolvesRingTopology(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.MapTopology("Unknown")
	playerLabels := []string{"A", "B"}
	lookup := test_helpers.NewTopologyServiceLookup(test_helpers.NewZoneFactories())

	// Act
	variant := lookup.Resolve(configuration.Topology)(
		*configuration, playerLabels, neutral_zone.Plans{},
		test_helpers.NewGenerationTuning(configuration, len(playerLabels)), "")

	// Assert: the ring closes back on itself, so two labels yield two connections.
	assert.Len(t, variant.Connections, 2)
}

func TestWhenTopologyIsRing_ResolvesRingTopology(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyRing
	playerLabels := []string{"A", "B"}
	lookup := test_helpers.NewTopologyServiceLookup(test_helpers.NewZoneFactories())

	// Act
	variant := lookup.Resolve(configuration.Topology)(
		*configuration, playerLabels, neutral_zone.Plans{},
		test_helpers.NewGenerationTuning(configuration, len(playerLabels)), "")

	// Assert
	assert.Len(t, variant.Connections, 2)
}

func TestWhenCityHoldIsEnabledForHubTopology_PassesHoldCityFlagToService(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyHubAndSpoke
	configuration.GameEndConditions = &config.GameEndConditions{CityHold: true}
	playerLabels := []string{"A", "B", "C"}
	lookup := test_helpers.NewTopologyServiceLookup(test_helpers.NewZoneFactories())

	// Act
	variant := lookup.Resolve(configuration.Topology)(
		*configuration, playerLabels, neutral_zone.Plans{},
		test_helpers.NewGenerationTuning(configuration, len(playerLabels)), "")

	// Assert
	assert.True(t, holdsCity(variant, "Hub"))
}

func hasZoneNamed(variant entities.Variant, name string) bool {
	for _, zone := range variant.Zones {
		if zone.Name == name {
			return true
		}
	}
	return false
}

func holdsCity(variant entities.Variant, zoneName string) bool {
	for _, zone := range variant.Zones {
		if zone.Name != zoneName {
			continue
		}
		for _, mainObject := range zone.MainObjects {
			if mainObject.HoldCityWinCon {
				return true
			}
		}
	}
	return false
}
