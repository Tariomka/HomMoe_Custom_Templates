package zoneLabelProvider_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
	"github.com/stretchr/testify/assert"
)

func mediumPlans(labels ...string) []models.NeutralZonePlan {
	plans := make([]models.NeutralZonePlan, len(labels))
	for i, label := range labels {
		plans[i] = models.NeutralZonePlan{Label: label, Quality: models.QualityMedium}
	}
	return plans
}

func orderingConfig(topology config.MapTopology, playerCount, neutralCount, separation int) config.GeneratorConfig {
	configuration := config.NewGeneratorConfig()
	configuration.Topology = topology
	configuration.PlayerCount = playerCount
	configuration.ZoneConfiguration.NeutralZoneCount = neutralCount
	configuration.MinNeutralZonesBetweenPlayers = separation
	return *configuration
}

func TestWhenSeparationIsZero_AppendsNeutralsAfterPlayers(t *testing.T) {
	t.Parallel()
	// Arrange
	provider := zones.NewZoneLabelProvider()
	configuration := orderingConfig(config.TopologyRing, 2, 2, 0)

	// Act
	ordered := provider.CreateOrderedZoneLabels(configuration, []string{"A", "B"}, mediumPlans("C", "D"), true)

	// Assert
	assert.Equal(t, []string{"A", "B", "C", "D"}, ordered)
}

func TestWhenRandomPortalsAreEnabled_IgnoresSeparation(t *testing.T) {
	t.Parallel()
	// Arrange
	provider := zones.NewZoneLabelProvider()
	configuration := orderingConfig(config.TopologyRing, 2, 2, 1)
	configuration.RandomPortals = true

	// Act
	ordered := provider.CreateOrderedZoneLabels(configuration, []string{"A", "B"}, mediumPlans("C", "D"), true)

	// Assert
	assert.Equal(t, []string{"A", "B", "C", "D"}, ordered)
}

func TestWhenSeparationCannotBeHonored_AppendsNeutralsAfterPlayers(t *testing.T) {
	t.Parallel()
	// Arrange
	provider := zones.NewZoneLabelProvider()
	configuration := orderingConfig(config.TopologyRing, 2, 1, 5)

	// Act
	ordered := provider.CreateOrderedZoneLabels(configuration, []string{"A", "B"}, mediumPlans("C"), true)

	// Assert
	assert.Equal(t, []string{"A", "B", "C"}, ordered)
}

func TestWhenRingSeparationIsHonored_InterleavesNeutralsBetweenPlayers(t *testing.T) {
	t.Parallel()
	// Arrange
	provider := zones.NewZoneLabelProvider()
	configuration := orderingConfig(config.TopologyRing, 2, 2, 1)

	// Act
	ordered := provider.CreateOrderedZoneLabels(configuration, []string{"A", "B"}, mediumPlans("C", "D"), true)

	// Assert
	assert.Equal(t, []string{"A", "C", "B", "D"}, ordered)
}

func TestWhenChainSeparationIsHonored_LeavesTrailingPlayerWithoutSeparator(t *testing.T) {
	t.Parallel()
	// Arrange
	provider := zones.NewZoneLabelProvider()
	configuration := orderingConfig(config.TopologyChain, 2, 1, 1)

	// Act
	ordered := provider.CreateOrderedZoneLabels(configuration, []string{"A", "B"}, mediumPlans("C"), false)

	// Assert
	assert.Equal(t, []string{"A", "C", "B"}, ordered)
}

func TestWhenMoreNeutralsThanSeparationSlots_AppendsLeftoversAtEnd(t *testing.T) {
	t.Parallel()
	// Arrange
	provider := zones.NewZoneLabelProvider()
	configuration := orderingConfig(config.TopologyRing, 2, 3, 1)

	// Act
	ordered := provider.CreateOrderedZoneLabels(configuration, []string{"A", "B"}, mediumPlans("C", "D", "E"), true)

	// Assert
	assert.Equal(t, []string{"A", "C", "B", "D", "E"}, ordered)
}

func TestWhenNoLabelsProvidedWithSeparation_ReturnsEmptySlice(t *testing.T) {
	t.Parallel()
	// Arrange
	provider := zones.NewZoneLabelProvider()
	configuration := orderingConfig(config.TopologyRing, 0, 0, 1)

	// Act
	ordered := provider.CreateOrderedZoneLabels(configuration, nil, nil, true)

	// Assert
	assert.Empty(t, ordered)
}

func TestWhenCirclesTopologyIsRing_DelegatesToBalancedRingOrdering(t *testing.T) {
	t.Parallel()
	// Arrange
	provider := zones.NewZoneLabelProvider()
	configuration := orderingConfig(config.TopologyCircles, 2, 2, 0)

	// Act
	ordered := provider.CreateOrderedZoneLabels(configuration, []string{"A", "B"}, mediumPlans("C", "D"), true)

	// Assert
	assert.Equal(t, []string{"A", "C", "B", "D"}, ordered)
}

func TestWhenCirclesTopologyIsChain_DelegatesToBalancedChainOrdering(t *testing.T) {
	t.Parallel()
	// Arrange
	provider := zones.NewZoneLabelProvider()
	configuration := orderingConfig(config.TopologyCircles, 1, 1, 0)

	// Act
	ordered := provider.CreateOrderedZoneLabels(configuration, []string{"A"}, mediumPlans("C"), false)

	// Assert
	assert.Equal(t, []string{"A", "C"}, ordered)
}

func TestWhenCirclesSeparationIsHonored_KeepsNeutralsBetweenPlayers(t *testing.T) {
	t.Parallel()
	// Arrange
	provider := zones.NewZoneLabelProvider()
	configuration := orderingConfig(config.TopologyCircles, 2, 2, 1)

	// Act
	ordered := provider.CreateOrderedZoneLabels(configuration, []string{"A", "B"}, mediumPlans("C", "D"), true)

	// Assert
	assert.Equal(t, []string{"A", "C", "B", "D"}, ordered)
}
