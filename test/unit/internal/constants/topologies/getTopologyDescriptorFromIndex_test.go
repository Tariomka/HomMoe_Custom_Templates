package topologies_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/common"
	"github.com/stretchr/testify/assert"
)

func TestWhenIndexIsWithinRange_ReturnsDescriptorAtIndex(t *testing.T) {
	t.Parallel()
	// Arrange
	expected := common.GetTopologyDescriptors().Circles

	// Act
	actual := common.GetTopologyDescriptorFromIndex(2)

	// Assert
	assert.Equal(t, expected, actual)
}

func TestWhenIndexIsNegative_ReturnsFirstTopology(t *testing.T) {
	t.Parallel()
	// Arrange
	expected := common.GetTopologyDescriptors().Random

	// Act
	actual := common.GetTopologyDescriptorFromIndex(-1)

	// Assert
	assert.Equal(t, expected, actual)
}

func TestWhenIndexIsBeyondRange_ReturnsFirstTopology(t *testing.T) {
	t.Parallel()
	// Arrange
	expected := common.GetTopologyDescriptors().Random

	// Act
	actual := common.GetTopologyDescriptorFromIndex(1000)

	// Assert
	assert.Equal(t, expected, actual)
}
