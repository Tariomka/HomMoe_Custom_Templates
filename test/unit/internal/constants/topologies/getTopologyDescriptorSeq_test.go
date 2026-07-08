package topologies_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/constants"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/stretchr/testify/assert"
)

func TestWhenSequenceIsIterated_YieldsAllTopologiesInDropdownOrder(t *testing.T) {
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
	for descriptor := range constants.GetTopologyDescriptorSeq() {
		actual = append(actual, descriptor.Type)
	}

	// Assert
	assert.Equal(t, expected, actual)
}

func TestWhenIterationStopsEarly_YieldsOnlyConsumedPrefix(t *testing.T) {
	// Arrange
	expected := []config.MapTopology{config.TopologyRandom, config.TopologyRing}
	var actual []config.MapTopology

	// Act
	for descriptor := range constants.GetTopologyDescriptorSeq() {
		actual = append(actual, descriptor.Type)
		if len(actual) == 2 {
			break
		}
	}

	// Assert
	assert.Equal(t, expected, actual)
}
