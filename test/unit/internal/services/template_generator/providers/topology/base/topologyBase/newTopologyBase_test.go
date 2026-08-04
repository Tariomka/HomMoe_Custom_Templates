package topologyBase_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/base"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
)

func TestWhenBaseIsConstructed_ProvidesZoneLabelProvider(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())

	// Assert
	assert.NotNil(t, topologyBase.ZoneLabelProvider)
}
