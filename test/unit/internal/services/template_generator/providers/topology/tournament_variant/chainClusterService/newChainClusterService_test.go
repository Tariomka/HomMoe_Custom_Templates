package chainClusterService_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/tournament_variant"
	"github.com/stretchr/testify/assert"
)

func TestWhenServiceIsConstructed_InitializesZoneLabelProvider(t *testing.T) {
	// Arrange & Act
	service := tournament_variant.NewChainClusterService()

	// Assert
	assert.NotNil(t, service.ZoneLabelProvider)
}
