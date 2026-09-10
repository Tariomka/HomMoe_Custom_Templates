package topologyConnectionService_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/base"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
	"github.com/stretchr/testify/assert"
)

func TestWhenBothLabelsArePlayers_ReturnsPlayerBorderGuardValue(t *testing.T) {
	t.Parallel()
	// Arrange
	connectionService := base.NewTopologyConnectionService(zones.NewZoneLabelProvider())

	// Act
	guardValue := connectionService.GetBorderGuardValue("A", "B", []string{"A", "B"}, nil, newUnitTuning())

	// Assert
	assert.Equal(t, 30000, guardValue)
}
