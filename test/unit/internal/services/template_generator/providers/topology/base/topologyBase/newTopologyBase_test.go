package topologyBase_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/base"
	"github.com/stretchr/testify/assert"
)

func TestWhenBaseIsConstructed_ProvidesZoneLabelProvider(t *testing.T) {
	// Arrange & Act
	topologyBase := base.NewTopologyBase()

	// Assert
	assert.NotNil(t, topologyBase.ZoneLabelProvider)
}
