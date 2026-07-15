package zoneLabelProvider_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutralZone"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
	"github.com/stretchr/testify/assert"
)

func simpleCountConfig(playerCount, neutralCount, castleCount int) config.GeneratorConfig {
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyRing
	configuration.PlayerCount = playerCount
	configuration.ZoneConfiguration.NeutralZoneCount = neutralCount
	configuration.ZoneConfiguration.NeutralZoneCastles = castleCount
	return *configuration
}

func TestWhenSimpleCountIsUsed_CreatesMediumPlansStartingAfterPlayerLabels(t *testing.T) {
	t.Parallel()
	// Arrange
	provider := zones.NewZoneLabelProvider()
	configuration := simpleCountConfig(2, 3, 2)
	expected := neutralZone.Plans{
		{Label: "C", Quality: neutralZone.QualityMedium, CastleCount: 2},
		{Label: "D", Quality: neutralZone.QualityMedium, CastleCount: 2},
		{Label: "E", Quality: neutralZone.QualityMedium, CastleCount: 2},
	}

	// Act
	plans := provider.CreateNeutralZonePlans(configuration)

	// Assert
	assert.Equal(t, expected, plans)
}

func TestWhenSimpleCastleCountExceedsFour_ClampsCastlesToFour(t *testing.T) {
	t.Parallel()
	// Arrange
	provider := zones.NewZoneLabelProvider()
	configuration := simpleCountConfig(2, 1, 9)

	// Act
	plans := provider.CreateNeutralZonePlans(configuration)

	// Assert
	assert.Equal(t, neutralZone.Plans{{Label: "C", Quality: neutralZone.QualityMedium, CastleCount: 4}}, plans)
}

func TestWhenAdvancedLowestCountsAreSet_CreatesLowestPlansBeforeLowPlans(t *testing.T) {
	t.Parallel()
	// Arrange
	provider := zones.NewZoneLabelProvider()
	configuration := simpleCountConfig(2, 0, 0)
	configuration.ZoneConfiguration.Advanced.NeutralLowestNoCastleCount = 1
	configuration.ZoneConfiguration.Advanced.NeutralLowestCastleCount = 1
	configuration.ZoneConfiguration.Advanced.NeutralLowestCastlesPerZone = 2
	configuration.ZoneConfiguration.Advanced.NeutralLowCastleCount = 1
	configuration.ZoneConfiguration.Advanced.NeutralLowCastlesPerZone = 1
	expected := neutralZone.Plans{
		{Label: "C", Quality: neutralZone.QualityLowest, CastleCount: 0},
		{Label: "D", Quality: neutralZone.QualityLowest, CastleCount: 2},
		{Label: "E", Quality: neutralZone.QualityLow, CastleCount: 1},
	}

	// Act
	plans := provider.CreateNeutralZonePlans(configuration)

	// Assert
	assert.Equal(t, expected, plans)
}

func TestWhenAdvancedCountsAreSet_CreatesPlansInTierOrder(t *testing.T) {
	t.Parallel()
	// Arrange
	provider := zones.NewZoneLabelProvider()
	configuration := simpleCountConfig(2, 0, 0)
	configuration.ZoneConfiguration.Advanced.NeutralLowNoCastleCount = 1
	configuration.ZoneConfiguration.Advanced.NeutralLowCastleCount = 1
	configuration.ZoneConfiguration.Advanced.NeutralLowCastlesPerZone = 2
	configuration.ZoneConfiguration.Advanced.NeutralMediumNoCastleCount = 1
	configuration.ZoneConfiguration.Advanced.NeutralMediumCastleCount = 1
	configuration.ZoneConfiguration.Advanced.NeutralMediumCastlesPerZone = 1
	configuration.ZoneConfiguration.Advanced.NeutralHighNoCastleCount = 1
	configuration.ZoneConfiguration.Advanced.NeutralHighCastleCount = 1
	configuration.ZoneConfiguration.Advanced.NeutralHighCastlesPerZone = 3
	expected := neutralZone.Plans{
		{Label: "C", Quality: neutralZone.QualityLow, CastleCount: 0},
		{Label: "D", Quality: neutralZone.QualityLow, CastleCount: 2},
		{Label: "E", Quality: neutralZone.QualityMedium, CastleCount: 0},
		{Label: "F", Quality: neutralZone.QualityMedium, CastleCount: 1},
		{Label: "G", Quality: neutralZone.QualityHigh, CastleCount: 0},
		{Label: "H", Quality: neutralZone.QualityHigh, CastleCount: 3},
	}

	// Act
	plans := provider.CreateNeutralZonePlans(configuration)

	// Assert
	assert.Equal(t, expected, plans)
}

func TestWhenAdvancedCastlesPerZoneExceedsFour_ClampsCastlesToFour(t *testing.T) {
	t.Parallel()
	// Arrange
	provider := zones.NewZoneLabelProvider()
	configuration := simpleCountConfig(2, 0, 0)
	configuration.ZoneConfiguration.Advanced.NeutralHighCastleCount = 1
	configuration.ZoneConfiguration.Advanced.NeutralHighCastlesPerZone = 7

	// Act
	plans := provider.CreateNeutralZonePlans(configuration)

	// Assert
	assert.Equal(t, neutralZone.Plans{{Label: "C", Quality: neutralZone.QualityHigh, CastleCount: 4}}, plans)
}

func TestWhenTopologyIsSharedWebAndNoNeutralsRequested_AddsSingleMediumPlan(t *testing.T) {
	t.Parallel()
	// Arrange
	provider := zones.NewZoneLabelProvider()
	configuration := simpleCountConfig(2, 0, 1)
	configuration.Topology = config.TopologySharedWeb

	// Act
	plans := provider.CreateNeutralZonePlans(configuration)

	// Assert
	assert.Equal(t, neutralZone.Plans{{Label: "C", Quality: neutralZone.QualityMedium, CastleCount: 1}}, plans)
}

func TestWhenRequestedCountExceedsLabelPool_CapsPlansAtAvailableLabels(t *testing.T) {
	t.Parallel()
	// Arrange
	provider := zones.NewZoneLabelProvider()
	configuration := simpleCountConfig(30, 5, 0)

	// Act
	plans := provider.CreateNeutralZonePlans(configuration)

	// Assert
	assert.Len(t, plans, 2)
}

func TestWhenRequestedCountIsNegative_CreatesNoPlans(t *testing.T) {
	t.Parallel()
	// Arrange
	provider := zones.NewZoneLabelProvider()
	configuration := simpleCountConfig(2, -3, 0)

	// Act
	plans := provider.CreateNeutralZonePlans(configuration)

	// Assert
	assert.Empty(t, plans)
}
