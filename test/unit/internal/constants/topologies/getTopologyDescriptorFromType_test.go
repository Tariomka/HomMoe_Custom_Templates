package topologies_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/constants"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/stretchr/testify/assert"
)

func TestWhenTypeIsKnown_ReturnsMatchingDescriptor(t *testing.T) {
	// Arrange
	expected := constants.GetTopologyDescriptors().Fractal

	// Act
	actual := constants.GetTopologyDescriptorFromType(config.TopologyFractal)

	// Assert
	assert.Equal(t, expected, actual)
}

func TestWhenTypeIsUnknown_ReturnsFirstTopology(t *testing.T) {
	// Arrange
	expected := constants.GetTopologyDescriptors().Random

	// Act
	actual := constants.GetTopologyDescriptorFromType(config.MapTopology("NoSuchTopology"))

	// Assert
	assert.Equal(t, expected, actual)
}
