package positionedTopologyBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
	"github.com/stretchr/testify/assert"
)

func TestWhenCreationServicesAreProvided_InitializesZoneLabelProvider(t *testing.T) {
	t.Parallel()
	// Arrange
	creationServices := zones.NewCreationServices(nil, nil)

	// Act
	builder := topology.NewPositionedTopologyBuilderWithCreationServices(creationServices)

	// Assert
	assert.NotNil(t, builder.ZoneLabelProvider)
}
