package squareTopology_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology"
	"github.com/stretchr/testify/assert"
)

func TestWhenServiceIsConstructed_InitializesZoneLabelProvider(t *testing.T) {
	// Arrange & Act
	service := topology.NewSquareTopologyService()

	// Assert
	assert.NotNil(t, service.ZoneLabelProvider)
}
