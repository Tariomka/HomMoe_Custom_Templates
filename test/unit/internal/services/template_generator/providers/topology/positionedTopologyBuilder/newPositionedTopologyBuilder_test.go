package positionedTopologyBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology"
	"github.com/stretchr/testify/assert"
)

func TestWhenCreated_InitializesZoneLabelProvider(t *testing.T) {
	t.Parallel()
	// Arrange

	// Act
	builder := topology.NewPositionedTopologyBuilder()

	// Assert
	assert.NotNil(t, builder.ZoneLabelProvider)
}
