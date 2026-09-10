package topologyConnectionService_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/base"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
	"github.com/stretchr/testify/assert"
)

func TestWhenTwoZonesAreRequested_OnePortalConnectionIsCreated(t *testing.T) {
	t.Parallel()
	// Arrange
	connectionService := base.NewTopologyConnectionService(zones.NewZoneLabelProvider())

	// Act
	connections := connectionService.CreateRandomPortalConnections(
		[]string{"A", "B"}, []string{"A", "B"}, newUnitTuning(), 1, nil)

	// Assert
	assert.Len(t, connections, 1)
}
