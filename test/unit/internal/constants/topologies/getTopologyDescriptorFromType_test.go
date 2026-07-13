package topologies_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/common"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/stretchr/testify/assert"
)

func TestWhenTypeIsKnown_ReturnsMatchingDescriptor(t *testing.T) {
	t.Parallel()
	// Arrange
	expected := common.GetTopologyDescriptors().Fractal

	// Act
	actual := common.GetTopologyDescriptorFromType(config.TopologyFractal)

	// Assert
	assert.Equal(t, expected, actual)
}

func TestWhenTypeIsUnknown_ReturnsFirstTopology(t *testing.T) {
	t.Parallel()
	// Arrange
	expected := common.GetTopologyDescriptors().Random

	// Act
	actual := common.GetTopologyDescriptorFromType(config.MapTopology("NoSuchTopology"))

	// Assert
	assert.Equal(t, expected, actual)
}
