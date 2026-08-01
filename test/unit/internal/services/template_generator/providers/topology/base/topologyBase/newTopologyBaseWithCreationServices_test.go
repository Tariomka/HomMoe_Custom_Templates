package topologyBase_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/base"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
	"github.com/stretchr/testify/assert"
)

func TestWhenCreationServicesAreProvided_ReturnsUsableTopologyBase(t *testing.T) {
	t.Parallel()
	// Arrange
	creationServices := zones.NewCreationServices(nil, nil)

	// Act
	topologyBase := base.NewTopologyBaseWithCreationServices(creationServices)

	// Assert
	assert.NotNil(t, topologyBase.ZoneLabelProvider)
}
