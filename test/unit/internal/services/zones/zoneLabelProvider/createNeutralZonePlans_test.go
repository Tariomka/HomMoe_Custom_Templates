package zoneLabelProvider_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
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
	// Arrange
	provider := zones.NewZoneLabelProvider()
	configuration := simpleCountConfig(2, 3, 2)
	expected := models.NeutralZonePlans{
		{Label: "C", Quality: models.QualityMedium, CastleCount: 2},
		{Label: "D", Quality: models.QualityMedium, CastleCount: 2},
		{Label: "E", Quality: models.QualityMedium, CastleCount: 2},
	}

	// Act
	plans := provider.CreateNeutralZonePlans(configuration)

	// Assert
	assert.Equal(t, expected, plans)
}

func TestWhenSimpleCastleCountExceedsFour_ClampsCastlesToFour(t *testing.T) {
	// Arrange
	provider := zones.NewZoneLabelProvider()
	configuration := simpleCountConfig(2, 1, 9)

	// Act
	plans := provider.CreateNeutralZonePlans(configuration)

	// Assert
	assert.Equal(t, models.NeutralZonePlans{{Label: "C", Quality: models.QualityMedium, CastleCount: 4}}, plans)
}

func TestWhenAdvancedCountsAreSet_CreatesPlansInTierOrder(t *testing.T) {
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
	expected := models.NeutralZonePlans{
		{Label: "C", Quality: models.QualityLow, CastleCount: 0},
		{Label: "D", Quality: models.QualityLow, CastleCount: 2},
		{Label: "E", Quality: models.QualityMedium, CastleCount: 0},
		{Label: "F", Quality: models.QualityMedium, CastleCount: 1},
		{Label: "G", Quality: models.QualityHigh, CastleCount: 0},
		{Label: "H", Quality: models.QualityHigh, CastleCount: 3},
	}

	// Act
	plans := provider.CreateNeutralZonePlans(configuration)

	// Assert
	assert.Equal(t, expected, plans)
}

func TestWhenAdvancedCastlesPerZoneExceedsFour_ClampsCastlesToFour(t *testing.T) {
	// Arrange
	provider := zones.NewZoneLabelProvider()
	configuration := simpleCountConfig(2, 0, 0)
	configuration.ZoneConfiguration.Advanced.NeutralHighCastleCount = 1
	configuration.ZoneConfiguration.Advanced.NeutralHighCastlesPerZone = 7

	// Act
	plans := provider.CreateNeutralZonePlans(configuration)

	// Assert
	assert.Equal(t, models.NeutralZonePlans{{Label: "C", Quality: models.QualityHigh, CastleCount: 4}}, plans)
}

func TestWhenTopologyIsSharedWebAndNoNeutralsRequested_AddsSingleMediumPlan(t *testing.T) {
	// Arrange
	provider := zones.NewZoneLabelProvider()
	configuration := simpleCountConfig(2, 0, 1)
	configuration.Topology = config.TopologySharedWeb

	// Act
	plans := provider.CreateNeutralZonePlans(configuration)

	// Assert
	assert.Equal(t, models.NeutralZonePlans{{Label: "C", Quality: models.QualityMedium, CastleCount: 1}}, plans)
}

func TestWhenRequestedCountExceedsLabelPool_CapsPlansAtAvailableLabels(t *testing.T) {
	// Arrange
	provider := zones.NewZoneLabelProvider()
	configuration := simpleCountConfig(30, 5, 0)

	// Act
	plans := provider.CreateNeutralZonePlans(configuration)

	// Assert
	assert.Len(t, plans, 2)
}

func TestWhenRequestedCountIsNegative_CreatesNoPlans(t *testing.T) {
	// Arrange
	provider := zones.NewZoneLabelProvider()
	configuration := simpleCountConfig(2, -3, 0)

	// Act
	plans := provider.CreateNeutralZonePlans(configuration)

	// Assert
	assert.Empty(t, plans)
}
