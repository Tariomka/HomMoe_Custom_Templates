package geometricTopology_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
	"github.com/stretchr/testify/assert"
)

func TestWhenCreationServicesAreProvided_InitializesService(t *testing.T) {
	t.Parallel()
	// Arrange
	creationServices := zones.NewCreationServices(nil, nil)

	// Act
	service := topology.NewGeometricTopologyServiceWithCreationServices(creationServices)

	// Assert
	assert.NotNil(t, service.ZoneLabelProvider)
}
