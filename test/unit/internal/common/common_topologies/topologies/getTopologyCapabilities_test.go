package topologies_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_topologies"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/stretchr/testify/assert"
)

func TestWhenCapabilitiesAreRequested_ReturnsRegisteredTopologyCapabilities(t *testing.T) {
	t.Parallel()
	// Arrange
	expected := map[config.MapTopology]models.TopologyCapabilities{
		config.TopologyRing:        {LayoutKind: models.TopologyLayoutRingHub},
		config.TopologyHubAndSpoke: {LayoutKind: models.TopologyLayoutRingHub, UsesHub: true},
		config.TopologyChain:       {LayoutKind: models.TopologyLayoutRingHub},
		config.TopologySharedWeb:   {LayoutKind: models.TopologyLayoutRingHub},
		config.TopologyRandom: {
			LayoutKind: models.TopologyLayoutScatter, UsesGeneratorPosition: true,
		},
		config.TopologyCircles: {
			LayoutKind: models.TopologyLayoutScatter, UsesGeneratorPosition: true, UsesGeneratorRing: true,
		},
		config.TopologySquare: {
			LayoutKind: models.TopologyLayoutFixedGeometry, UsesGeneratorPosition: true,
		},
		config.TopologyGeometric: {
			LayoutKind: models.TopologyLayoutFixedGeometry, UsesGeneratorPosition: true,
		},
		config.TopologyCross: {
			LayoutKind: models.TopologyLayoutFixedGeometry, UsesGeneratorPosition: true,
		},
		config.TopologyFractal: {
			LayoutKind: models.TopologyLayoutFixedGeometry, UsesGeneratorPosition: true,
		},
		config.TopologyGeometricHub: {
			LayoutKind: models.TopologyLayoutFixedGeometry, UsesHub: true, UsesGeneratorPosition: true,
		},
		config.MapTopology("Unknown"): {LayoutKind: models.TopologyLayoutRingHub},
	}
	actual := make(map[config.MapTopology]models.TopologyCapabilities, len(expected))

	// Act
	for topology := range expected {
		actual[topology] = common_topologies.GetTopologyCapabilities(topology)
	}

	// Assert
	assert.Equal(t, expected, actual)
}
