package topologyBase_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/base"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
)

func TestWhenFirstLabelIsPlayer_ZeroAngleZoneIsThatSpawnZone(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())

	// Act
	variant := topologyBase.CreateVariant([]string{"A", "B"}, "A", 4, nil, nil)

	// Assert
	assert.Equal(t, "Spawn-A", variant.Orientation.ZeroAngleZone)
}

func TestWhenFirstLabelIsNeutral_ZeroAngleZoneIsThatNeutralZone(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())

	// Act
	variant := topologyBase.CreateVariant([]string{"A", "B"}, "C", 4, nil, nil)

	// Assert
	assert.Equal(t, "Neutral-C", variant.Orientation.ZeroAngleZone)
}

func TestWhenZoneCountIsPositive_RandomAngleStepDividesFullCircleByZoneCount(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())

	// Act
	variant := topologyBase.CreateVariant([]string{"A"}, "A", 8, nil, nil)

	// Assert
	assert.Equal(t, 45, variant.Orientation.RandomAngleStep)
}

func TestWhenZoneCountIsZero_RandomAngleStepStaysUnset(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())

	// Act
	variant := topologyBase.CreateVariant([]string{"A"}, "A", 0, nil, nil)

	// Assert
	assert.Equal(t, 0, variant.Orientation.RandomAngleStep)
}

func TestWhenZonesProvided_VariantCarriesZonesVerbatim(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())
	zones := []template_model.Zone{{Name: "Spawn-A"}, {Name: "Neutral-C"}}

	// Act
	variant := topologyBase.CreateVariant([]string{"A"}, "A", 2, zones, nil)

	// Assert
	assert.Equal(t, zones, variant.Zones)
}

func TestWhenConnectionsProvided_VariantCarriesConnectionsVerbatim(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())
	connections := []template_model.Connection{{Name: "Ring-A-C", From: "Spawn-A", To: "Neutral-C"}}

	// Act
	variant := topologyBase.CreateVariant([]string{"A"}, "A", 2, nil, connections)

	// Assert
	assert.Equal(t, connections, variant.Connections)
}
