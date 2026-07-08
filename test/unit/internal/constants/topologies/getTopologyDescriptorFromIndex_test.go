package topologies_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/constants"
	"github.com/stretchr/testify/assert"
)

func TestWhenIndexIsWithinRange_ReturnsDescriptorAtIndex(t *testing.T) {
	// Arrange
	expected := constants.GetTopologyDescriptors().Circles

	// Act
	actual := constants.GetTopologyDescriptorFromIndex(2)

	// Assert
	assert.Equal(t, expected, actual)
}

func TestWhenIndexIsNegative_ReturnsFirstTopology(t *testing.T) {
	// Arrange
	expected := constants.GetTopologyDescriptors().Random

	// Act
	actual := constants.GetTopologyDescriptorFromIndex(-1)

	// Assert
	assert.Equal(t, expected, actual)
}

func TestWhenIndexIsBeyondRange_ReturnsFirstTopology(t *testing.T) {
	// Arrange
	expected := constants.GetTopologyDescriptors().Random

	// Act
	actual := constants.GetTopologyDescriptorFromIndex(1000)

	// Assert
	assert.Equal(t, expected, actual)
}
