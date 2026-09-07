package zoneLabelProvider_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenNoNeutralPlansExist_ReturnsEmptyLabel(t *testing.T) {
	t.Parallel()
	// Arrange
	provider := zones.NewZoneLabelProvider()
	configuration := hubCityHoldConfig(gofakeit.Number(2, 8))

	// Act
	label := provider.GetHoldCityLabel(configuration, []string{"A", "B"}, neutral_zone.Plans{})

	// Assert
	assert.Empty(t, label)
}

func TestWhenTopologyIsNotHubAndSpoke_ReturnsEmptyLabel(t *testing.T) {
	t.Parallel()
	// Arrange
	provider := zones.NewZoneLabelProvider()
	configuration := hubCityHoldConfig(gofakeit.Number(2, 8))
	configuration.Topology = config.TopologyRing
	plans := neutral_zone.Plans{{Label: "C", Quality: neutral_zone.QualityMedium, CastleCount: 1}}

	// Act
	label := provider.GetHoldCityLabel(configuration, []string{"A", "B"}, plans)

	// Assert
	assert.Empty(t, label)
}

func TestWhenCityHoldModeIsOff_ReturnsEmptyLabel(t *testing.T) {
	t.Parallel()
	// Arrange
	provider := zones.NewZoneLabelProvider()
	configuration := hubCityHoldConfig(gofakeit.Number(2, 8))
	configuration.GameEndConditions.CityHold = false
	plans := neutral_zone.Plans{{Label: "C", Quality: neutral_zone.QualityMedium, CastleCount: 1}}

	// Act
	label := provider.GetHoldCityLabel(configuration, []string{"A", "B"}, plans)

	// Assert
	assert.Empty(t, label)
}

func TestWhenDistancesAndVariancesTie_PicksHigherQualityNeutral(t *testing.T) {
	t.Parallel()
	// Arrange
	provider := zones.NewZoneLabelProvider()
	configuration := hubCityHoldConfig(gofakeit.Number(2, 8))
	plans := neutral_zone.Plans{
		{Label: "C", Quality: neutral_zone.QualityMedium, CastleCount: 1},
		{Label: "D", Quality: neutral_zone.QualityLow, CastleCount: 0},
	}

	// Act
	label := provider.GetHoldCityLabel(configuration, []string{"A", "B"}, plans)

	// Assert
	assert.Equal(t, "C", label)
}

func TestWhenQualityAlsoTies_PicksNeutralWithCastle(t *testing.T) {
	t.Parallel()
	// Arrange
	provider := zones.NewZoneLabelProvider()
	configuration := hubCityHoldConfig(gofakeit.Number(2, 8))
	plans := neutral_zone.Plans{
		{Label: "C", Quality: neutral_zone.QualityMedium, CastleCount: 0},
		{Label: "D", Quality: neutral_zone.QualityMedium, CastleCount: 1},
	}

	// Act
	label := provider.GetHoldCityLabel(configuration, []string{"A", "B"}, plans)

	// Assert
	assert.Equal(t, "D", label)
}

func hubCityHoldConfig(playerCount int) config.GeneratorConfig {
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyHubAndSpoke
	configuration.PlayerCount = playerCount
	configuration.GameEndConditions.CityHold = true
	return *configuration
}
