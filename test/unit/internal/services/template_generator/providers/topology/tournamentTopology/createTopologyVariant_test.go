package tournamentTopology_test

import (
	"strings"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
)

func newChainTournamentConfig() *config.GeneratorConfig {
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyChain
	return configuration
}

func newFourNeutralPlans() neutral_zone.Plans {
	neutralZones := neutral_zone.Plans{}
	neutralZones.AddPlan("C", neutral_zone.QualityLow, 0)
	neutralZones.AddPlan("D", neutral_zone.QualityMedium, 1)
	neutralZones.AddPlan("E", neutral_zone.QualityMedium, 1)
	neutralZones.AddPlan("F", neutral_zone.QualityHigh, 1)
	return neutralZones
}

func TestWhenFourNeutralPlansAreSplitAcrossTwoPlayers_CreatesZonePerPlayerAndNeutralLabel(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := newChainTournamentConfig()
	neutralZones := newFourNeutralPlans()
	tuning := test_helpers.NewGenerationTuning(configuration, 6)
	service := topology.NewTournamentTopologyService(test_helpers.NewZoneFactories())

	// Act
	variant := service.CreateTopologyVariant(*configuration, []string{"A", "B"}, neutralZones, tuning)

	// Assert
	assert.Len(t, variant.Zones, 6)
}

func TestWhenTournamentIsBuilt_EveryConnectionReferencesExistingZones(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := newChainTournamentConfig()
	neutralZones := newFourNeutralPlans()
	tuning := test_helpers.NewGenerationTuning(configuration, 6)
	service := topology.NewTournamentTopologyService(test_helpers.NewZoneFactories())

	// Act
	variant := service.CreateTopologyVariant(*configuration, []string{"A", "B"}, neutralZones, tuning)

	// Assert
	assert.Empty(t, danglingConnectionNames(variant))
}

func TestWhenPortalsAreDisabled_PlayerClustersStayIsolatedAsTwoComponents(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := newChainTournamentConfig()
	configuration.RandomPortals = false
	neutralZones := newFourNeutralPlans()
	tuning := test_helpers.NewGenerationTuning(configuration, 6)
	service := topology.NewTournamentTopologyService(test_helpers.NewZoneFactories())

	// Act
	variant := service.CreateTopologyVariant(*configuration, []string{"A", "B"}, neutralZones, tuning)

	// Assert
	assert.Equal(t, 2, connectedComponentCount(variant))
}

func TestWhenTwoPlayersAreProvided_EachPlayerGetsOwnSpawnZone(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := newChainTournamentConfig()
	neutralZones := newFourNeutralPlans()
	tuning := test_helpers.NewGenerationTuning(configuration, 6)
	service := topology.NewTournamentTopologyService(test_helpers.NewZoneFactories())

	// Act
	variant := service.CreateTopologyVariant(*configuration, []string{"A", "B"}, neutralZones, tuning)

	// Assert
	names := zoneNameSet(variant)
	assert.True(t, names["Spawn-A"] && names["Spawn-B"],
		"expected both Spawn-A and Spawn-B in %v", names)
}

func TestWhenTopologyIsHubAndSpoke_CreatesHubZonePerPlayer(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyHubAndSpoke
	neutralZones := newFourNeutralPlans()
	tuning := test_helpers.NewGenerationTuning(configuration, 8)
	service := topology.NewTournamentTopologyService(test_helpers.NewZoneFactories())

	// Act
	variant := service.CreateTopologyVariant(*configuration, []string{"A", "B"}, neutralZones, tuning)

	// Assert
	names := zoneNameSet(variant)
	assert.True(t, names["Hub-A"] && names["Hub-B"],
		"expected both Hub-A and Hub-B in %v", names)
}

func TestWhenTopologyIsRing_CreatesRingClusterConnections(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyRing
	neutralZones := newFourNeutralPlans()
	tuning := test_helpers.NewGenerationTuning(configuration, 6)
	service := topology.NewTournamentTopologyService(test_helpers.NewZoneFactories())

	// Act
	variant := service.CreateTopologyVariant(*configuration, []string{"A", "B"}, neutralZones, tuning)

	// Assert
	assert.NotZero(t, countConnectionsWithPrefix(variant, "TRing-"))
}

func TestWhenTopologyIsCircles_CreatesBalancedClusterConnections(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyCircles
	neutralZones := newFourNeutralPlans()
	tuning := test_helpers.NewGenerationTuning(configuration, 6)
	service := topology.NewTournamentTopologyService(test_helpers.NewZoneFactories())

	// Act
	variant := service.CreateTopologyVariant(*configuration, []string{"A", "B"}, neutralZones, tuning)

	// Assert
	assert.NotZero(t, countConnectionsWithPrefix(variant, "TBal-"))
}

func TestWhenTopologyIsUnhandled_FallsBackToChainClusterConnections(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologySharedWeb
	neutralZones := newFourNeutralPlans()
	tuning := test_helpers.NewGenerationTuning(configuration, 6)
	service := topology.NewTournamentTopologyService(test_helpers.NewZoneFactories())

	// Act
	variant := service.CreateTopologyVariant(*configuration, []string{"A", "B"}, neutralZones, tuning)

	// Assert
	assert.NotZero(t, countConnectionsWithPrefix(variant, "Tourney-"))
}

func TestWhenRandomPortalsAreEnabled_AddsPortalConnections(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := newChainTournamentConfig()
	configuration.RandomPortals = true
	neutralZones := newFourNeutralPlans()
	tuning := test_helpers.NewGenerationTuning(configuration, 6)
	service := topology.NewTournamentTopologyService(test_helpers.NewZoneFactories())

	// Act
	variant := service.CreateTopologyVariant(*configuration, []string{"A", "B"}, neutralZones, tuning)

	// Assert
	assert.NotZero(t, countConnectionsWithPrefix(variant, "Portal-"))
}

func TestWhenNeutralPlansAreSplit_EachClusterGetsHalfOfNeutralZones(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := newChainTournamentConfig()
	configuration.RandomPortals = false
	neutralZones := newFourNeutralPlans()
	tuning := test_helpers.NewGenerationTuning(configuration, 6)
	service := topology.NewTournamentTopologyService(test_helpers.NewZoneFactories())

	// Act
	variant := service.CreateTopologyVariant(*configuration, []string{"A", "B"}, neutralZones, tuning)

	// Assert
	neutralPerCluster := map[string]int{}
	for index, zone := range variant.Zones {
		if strings.HasPrefix(zone.Name, "Neutral-") {
			// Chain clusters are emitted player-by-player: zones 0..2 belong to
			// player A, zones 3..5 to player B.
			clusterKey := "A"
			if index >= len(variant.Zones)/2 {
				clusterKey = "B"
			}
			neutralPerCluster[clusterKey]++
		}
	}
	assert.Equal(t, map[string]int{"A": 2, "B": 2}, neutralPerCluster)
}
