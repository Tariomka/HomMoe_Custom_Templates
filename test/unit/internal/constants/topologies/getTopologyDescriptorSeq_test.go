package topologies_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/common"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/stretchr/testify/assert"
)

func TestWhenSequenceIsIterated_YieldsAllTopologiesInDropdownOrder(t *testing.T) {
	t.Parallel()
	// Arrange
	expected := []config.MapTopology{
		config.TopologyRandom,
		config.TopologyRing,
		config.TopologyCircles,
		config.TopologyHubAndSpoke,
		config.TopologyChain,
		config.TopologySharedWeb,
		config.TopologySquare,
		config.TopologyGeometric,
		config.TopologyCross,
		config.TopologyFractal,
	}
	var actual []config.MapTopology

	// Act
	for descriptor := range common.GetTopologyDescriptorSeq() {
		actual = append(actual, descriptor.Type)
	}

	// Assert
	assert.Equal(t, expected, actual)
}

func TestWhenIterationStopsEarly_YieldsOnlyConsumedPrefix(t *testing.T) {
	t.Parallel()
	// Arrange
	expected := []config.MapTopology{config.TopologyRandom, config.TopologyRing}
	var actual []config.MapTopology

	// Act
	for descriptor := range common.GetTopologyDescriptorSeq() {
		actual = append(actual, descriptor.Type)
		if len(actual) == 2 {
			break
		}
	}

	// Assert
	assert.Equal(t, expected, actual)
}
