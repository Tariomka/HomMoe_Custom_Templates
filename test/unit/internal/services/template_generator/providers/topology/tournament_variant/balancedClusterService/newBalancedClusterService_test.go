package balancedClusterService_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/tournament_variant"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
)

func TestWhenServiceIsConstructed_InitializesZoneLabelProvider(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	service := tournament_variant.NewBalancedClusterService(test_helpers.NewZoneFactories())

	// Assert
	assert.NotNil(t, service.ZoneLabelProvider)
}
