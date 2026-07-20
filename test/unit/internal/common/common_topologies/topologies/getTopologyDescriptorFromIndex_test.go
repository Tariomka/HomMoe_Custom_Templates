package topologies_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_topologies"
	"github.com/stretchr/testify/assert"
)

func TestWhenIndexIsWithinRange_ReturnsDescriptorAtIndex(t *testing.T) {
	t.Parallel()
	// Arrange
	expected := common_topologies.GetTopologyDescriptors().Circles

	// Act
	actual := common_topologies.GetTopologyDescriptorFromIndex(2)

	// Assert
	assert.Equal(t, expected, actual)
}

func TestWhenIndexIsNegative_ReturnsFirstTopology(t *testing.T) {
	t.Parallel()
	// Arrange
	expected := common_topologies.GetTopologyDescriptors().Random

	// Act
	actual := common_topologies.GetTopologyDescriptorFromIndex(-1)

	// Assert
	assert.Equal(t, expected, actual)
}

func TestWhenIndexIsBeyondRange_ReturnsFirstTopology(t *testing.T) {
	t.Parallel()
	// Arrange
	expected := common_topologies.GetTopologyDescriptors().Random

	// Act
	actual := common_topologies.GetTopologyDescriptorFromIndex(1000)

	// Assert
	assert.Equal(t, expected, actual)
}
