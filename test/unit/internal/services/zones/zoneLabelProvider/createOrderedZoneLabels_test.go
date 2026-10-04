package zoneLabelProvider_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
	"github.com/stretchr/testify/assert"
)

func TestWhenTopologyIsNotCircles_AppendsNeutralsAfterPlayers(t *testing.T) {
	t.Parallel()
	// Arrange
	provider := zones.NewZoneLabelProvider()
	configuration := orderingConfig(config.TopologySquare, 2, 2)

	// Act
	ordered := provider.CreateOrderedZoneLabels(configuration, []string{"A", "B"}, mediumPlans("C", "D"))

	// Assert
	assert.Equal(t, []string{"A", "B", "C", "D"}, ordered)
}

func TestWhenNoLabelsProvided_ReturnsEmptySlice(t *testing.T) {
	t.Parallel()
	// Arrange
	provider := zones.NewZoneLabelProvider()
	configuration := orderingConfig(config.TopologySquare, 0, 0)

	// Act
	ordered := provider.CreateOrderedZoneLabels(configuration, nil, nil)

	// Assert
	assert.Empty(t, ordered)
}

func TestWhenTopologyIsCircles_DelegatesToBalancedRingOrdering(t *testing.T) {
	t.Parallel()
	// Arrange
	provider := zones.NewZoneLabelProvider()
	configuration := orderingConfig(config.TopologyCircles, 2, 2)

	// Act
	ordered := provider.CreateOrderedZoneLabels(configuration, []string{"A", "B"}, mediumPlans("C", "D"))

	// Assert
	assert.Equal(t, []string{"A", "C", "B", "D"}, ordered)
}

func mediumPlans(labels ...string) []neutral_zone.Plan {
	plans := make([]neutral_zone.Plan, len(labels))
	for i, label := range labels {
		plans[i] = neutral_zone.Plan{Label: label, Quality: neutral_zone.QualityMedium}
	}
	return plans
}

func orderingConfig(topology config.MapTopology, playerCount, neutralCount int) config.GeneratorConfig {
	configuration := config.NewGeneratorConfig()
	configuration.Topology = topology
	configuration.PlayerCount = playerCount
	configuration.ZoneConfiguration.NeutralZoneCount = neutralCount
	return *configuration
}
