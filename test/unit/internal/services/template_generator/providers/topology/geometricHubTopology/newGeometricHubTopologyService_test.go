package geometricHubTopology_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology"
	"github.com/stretchr/testify/assert"
)

func TestWhenServiceIsConstructed_ReturnsNonNilService(t *testing.T) {
	t.Parallel()
	// Arrange

	// Act
	service := topology.NewGeometricHubTopologyService()

	// Assert
	assert.NotNil(t, service)
}
