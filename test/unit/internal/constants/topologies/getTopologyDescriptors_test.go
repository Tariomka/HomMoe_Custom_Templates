package topologies_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/constants"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/stretchr/testify/assert"
)

func TestWhenCatalogIsRequested_EveryDescriptorCarriesItsOwnType(t *testing.T) {
	t.Parallel()
	// Arrange
	expected := []config.MapTopology{
		config.TopologyRing,
		config.TopologyCircles,
		config.TopologyRandom,
		config.TopologyHubAndSpoke,
		config.TopologyChain,
		config.TopologySharedWeb,
		config.TopologySquare,
		config.TopologyGeometric,
		config.TopologyCross,
		config.TopologyFractal,
	}

	// Act
	descriptors := constants.GetTopologyDescriptors()

	// Assert
	actual := []config.MapTopology{
		descriptors.Default.Type,
		descriptors.Circles.Type,
		descriptors.Random.Type,
		descriptors.HubAndSpoke.Type,
		descriptors.Chain.Type,
		descriptors.SharedWeb.Type,
		descriptors.Square.Type,
		descriptors.Geometric.Type,
		descriptors.Cross.Type,
		descriptors.Fractal.Type,
	}
	assert.Equal(t, expected, actual)
}
