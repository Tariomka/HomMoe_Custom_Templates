package hubTopology_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology"
	"github.com/stretchr/testify/assert"
)

func TestWhenServiceIsConstructed_InitializesZoneLabelProvider(t *testing.T) {
	// Arrange & Act
	service := topology.NewHubTopologyService()

	// Assert
	assert.NotNil(t, service.ZoneLabelProvider)
}
