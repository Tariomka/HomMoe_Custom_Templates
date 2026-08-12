package topologies_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_topologies"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/stretchr/testify/assert"
)

func TestWhenTypeIsKnown_ReturnsMatchingDescriptor(t *testing.T) {
	t.Parallel()
	// Arrange
	expected := common_topologies.GetTopologyDescriptors().Fractal

	// Act
	actual := common_topologies.GetTopologyDescriptorFromType(config.TopologyFractal)

	// Assert
	assert.Equal(t, expected, actual)
}

func TestWhenTypeIsUnknown_ReturnsRingFallback(t *testing.T) {
	t.Parallel()
	// Arrange
	expected := common_topologies.GetTopologyDescriptors().Default

	// Act
	actual := common_topologies.GetTopologyDescriptorFromType(config.MapTopology("NoSuchTopology"))

	// Assert
	assert.Equal(t, expected, actual)
}
