package topologies_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_topologies"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/stretchr/testify/assert"
)

func TestWhenCatalogIsRequested_EveryDescriptorCarriesItsOwnType(t *testing.T) {
	t.Parallel()
	// Arrange
	expected := []config.MapTopology{
		config.TopologyCircles,
		config.TopologyRandom,
		config.TopologyGeometricHub,
		config.TopologySquare,
		config.TopologyGeometric,
		config.TopologyCross,
		config.TopologyFractal,
	}

	// Act
	descriptors := common_topologies.GetTopologyDescriptors()

	// Assert
	actual := []config.MapTopology{
		descriptors.Circles.Type,
		descriptors.Random.Type,
		descriptors.GeometricHub.Type,
		descriptors.Square.Type,
		descriptors.Geometric.Type,
		descriptors.Cross.Type,
		descriptors.Fractal.Type,
	}
	assert.Equal(t, expected, actual)
}
