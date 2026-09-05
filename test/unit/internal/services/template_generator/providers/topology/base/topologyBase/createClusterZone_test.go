package topologyBase_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/base"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
)

func TestWhenClusterZoneIsSpawn_CreatesSpawnZoneForPlayerIndex(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())

	// Act
	zone := topologyBase.CreateClusterZone(
		newClusterConfiguration(), "A", nil, 2, true, false, newUnitTuning(), nil)

	// Assert
	assert.Equal(t, "Player3", zone.MainObjects[0].Spawn)
}

func TestWhenClusterZoneIsSpawn_NameCombinesSpawnPrefixWithLabel(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())

	// Act
	zone := topologyBase.CreateClusterZone(
		newClusterConfiguration(), "A", nil, 0, true, false, newUnitTuning(), nil)

	// Assert
	assert.Equal(t, "Spawn-A", zone.Name)
}

func TestWhenClusterZoneIsNeutral_NameCombinesNeutralPrefixWithMatchingPlanLabel(t *testing.T) {
	t.Parallel()
	// Arrange
	plans := neutral_zone.Plans{
		{Label: "B", Quality: neutral_zone.QualityMedium, CastleCount: 1},
		{Label: "C", Quality: neutral_zone.QualityLow, CastleCount: 0},
	}
	topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())

	// Act
	zone := topologyBase.CreateClusterZone(
		newClusterConfiguration(), "C", nil, 0, false, false, newUnitTuning(), plans)

	// Assert
	assert.Equal(t, "Neutral-C", zone.Name)
}

func TestWhenClusterZoneIsNeutralAndHoldsCity_PrimaryCastleCarriesHoldCityWinCondition(t *testing.T) {
	t.Parallel()
	// Arrange
	plans := neutral_zone.Plans{{Label: "B", Quality: neutral_zone.QualityHigh, CastleCount: 1}}
	topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())

	// Act
	zone := topologyBase.CreateClusterZone(
		newClusterConfiguration(), "B", nil, 0, false, true, newUnitTuning(), plans)

	// Assert
	assert.True(t, zone.MainObjects[0].HoldCityWinCon)
}

// newClusterConfiguration builds the minimal generator configuration that
// CreateClusterZone reads from.
func newClusterConfiguration() config.GeneratorConfig {
	configuration := config.GeneratorConfig{GenerateRoads: true}
	configuration.ZoneConfiguration.PlayerZoneCastles = 0
	configuration.ZoneConfiguration.PlayerZoneSize = 1.0
	configuration.ZoneConfiguration.NeutralZoneSize = 1.0

	return configuration
}
