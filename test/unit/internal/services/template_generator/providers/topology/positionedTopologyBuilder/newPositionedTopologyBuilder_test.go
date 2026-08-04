package positionedTopologyBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
)

func TestWhenCreated_InitializesZoneLabelProvider(t *testing.T) {
	t.Parallel()
	// Arrange

	// Act
	builder := topology.NewPositionedTopologyBuilder(test_helpers.NewZoneFactories())

	// Assert
	assert.NotNil(t, builder.ZoneLabelProvider)
}
