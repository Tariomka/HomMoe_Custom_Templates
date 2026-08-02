package chainTopology_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
	"github.com/stretchr/testify/assert"
)

func TestWhenServiceIsConstructed_InitializesZoneLabelProvider(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	service := topology.NewChainTopologyService(zones.NewCreationServices(nil, nil))

	// Assert
	assert.NotNil(t, service.ZoneLabelProvider)
}
