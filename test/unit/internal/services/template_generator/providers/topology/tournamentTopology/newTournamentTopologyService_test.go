package tournamentTopology_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
)

func TestWhenServiceIsConstructed_InitializesZoneLabelProvider(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	service := topology.NewTournamentTopologyService(test_helpers.NewTournamentTopologyDependencies())

	// Assert
	assert.NotNil(t, service.ZoneLabelProvider)
}
