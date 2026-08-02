package positionedTopologyBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
	"github.com/stretchr/testify/assert"
)

func TestWhenCreated_InitializesZoneLabelProvider(t *testing.T) {
	t.Parallel()
	// Arrange

	// Act
	builder := topology.NewPositionedTopologyBuilder(zones.NewCreationServices(nil, nil))

	// Assert
	assert.NotNil(t, builder.ZoneLabelProvider)
}
