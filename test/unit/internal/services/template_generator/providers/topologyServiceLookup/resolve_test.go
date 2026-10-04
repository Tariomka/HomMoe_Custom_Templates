package topologyServiceLookup_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_topologies"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// unsupportedTopologies are IDs no service exists for: the four retired ones,
// an unknown one and the empty one.
var unsupportedTopologies = []config.MapTopology{ //nolint:gochecknoglobals // Read-only test table.
	"Default", "HubAndSpoke", "Chain", "SharedWeb", "NotARealTopology", "",
}

func TestWhenTopologyIsSupported_ReportsFound(t *testing.T) {
	t.Parallel()
	for descriptor := range common_topologies.GetTopologyDescriptorSeq() {
		t.Run(string(descriptor.Type)+"_ReportsFound", func(t *testing.T) {
			t.Parallel()
			// Arrange
			lookup := test_helpers.NewTopologyServiceLookup(test_helpers.NewZoneFactories())

			// Act
			_, found := lookup.Resolve(descriptor.Type)

			// Assert
			assert.True(t, found)
		})
	}
}

func TestWhenTopologyIsSupported_ReturnsCreator(t *testing.T) {
	t.Parallel()
	for descriptor := range common_topologies.GetTopologyDescriptorSeq() {
		t.Run(string(descriptor.Type)+"_ReturnsCreator", func(t *testing.T) {
			t.Parallel()
			// Arrange
			lookup := test_helpers.NewTopologyServiceLookup(test_helpers.NewZoneFactories())

			// Act
			creator, _ := lookup.Resolve(descriptor.Type)

			// Assert
			assert.NotNil(t, creator)
		})
	}
}

func TestWhenTopologyIsNotSupported_ReportsNotFound(t *testing.T) {
	t.Parallel()
	for _, mapTopology := range unsupportedTopologies {
		t.Run("Topology"+string(mapTopology)+"_ReportsNotFound", func(t *testing.T) {
			t.Parallel()
			// Arrange
			lookup := test_helpers.NewTopologyServiceLookup(test_helpers.NewZoneFactories())

			// Act
			_, found := lookup.Resolve(mapTopology)

			// Assert
			assert.False(t, found)
		})
	}
}

func TestWhenTopologyIsNotSupported_ReturnsNoCreator(t *testing.T) {
	t.Parallel()
	for _, mapTopology := range unsupportedTopologies {
		t.Run("Topology"+string(mapTopology)+"_ReturnsNoCreator", func(t *testing.T) {
			t.Parallel()
			// Arrange
			lookup := test_helpers.NewTopologyServiceLookup(test_helpers.NewZoneFactories())

			// Act
			creator, _ := lookup.Resolve(mapTopology)

			// Assert
			assert.Nil(t, creator)
		})
	}
}

func TestWhenTopologyIsGeometricHub_ResolvesHubTopology(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyGeometricHub
	playerLabels := []string{"A", "B", "C"}
	lookup := test_helpers.NewTopologyServiceLookup(test_helpers.NewZoneFactories())
	creator, found := lookup.Resolve(configuration.Topology)
	require.True(t, found)

	// Act
	variant := creator(*configuration, playerLabels, neutral_zone.Plans{},
		test_helpers.NewGenerationTuning(configuration, len(playerLabels)), "")

	// Assert
	assert.True(t, hasZoneNamed(variant, "Hub"))
}

func TestWhenCityHoldIsEnabledForGeometricHub_PassesHoldCityFlagToService(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyGeometricHub
	configuration.GameEndConditions = &config.GameEndConditions{CityHold: true}
	playerLabels := []string{"A", "B", "C"}
	lookup := test_helpers.NewTopologyServiceLookup(test_helpers.NewZoneFactories())
	creator, found := lookup.Resolve(configuration.Topology)
	require.True(t, found)

	// Act
	variant := creator(*configuration, playerLabels, neutral_zone.Plans{},
		test_helpers.NewGenerationTuning(configuration, len(playerLabels)), "")

	// Assert
	assert.True(t, holdsCity(variant, "Hub"))
}

func hasZoneNamed(variant template_model.Variant, name string) bool {
	for _, zone := range variant.Zones {
		if zone.Name == name {
			return true
		}
	}
	return false
}

func holdsCity(variant template_model.Variant, zoneName string) bool {
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
