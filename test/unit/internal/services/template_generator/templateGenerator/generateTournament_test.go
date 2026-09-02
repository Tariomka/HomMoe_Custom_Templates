package templateGenerator_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/linq"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenTournamentEnabled_CreatesSpawnZonePerPlayer(t *testing.T) {
	t.Parallel()
	// Arrange
	playerCount := gofakeit.Number(2, 8)
	generator := test_helpers.NewTemplateGenerator(
		newTournamentConfiguration(config.TopologyRing, playerCount, gofakeit.Number(1, 20)))

	// Act
	actual, _ := generateTemplate(generator)

	// Assert
	assert.Len(t, zonesWithPrefix(actual, "Spawn-"), playerCount)
}

func TestWhenTournamentEnabledWithHubAndSpokeTopology_CreatesHubGuardGroups(t *testing.T) {
	t.Parallel()
	// Arrange
	generator := test_helpers.NewTemplateGenerator(
		newTournamentConfiguration(config.TopologyHubAndSpoke, 2, gofakeit.Number(1, 20)))

	// Act
	actual, _ := generateTemplate(generator)

	// Assert
	hasHubGuardGroup := linq.FromSlice(actual.Variants[0].Connections).
		Where(func(connection entities.Connection) bool {
			return strings.HasPrefix(connection.GuardMatchGroup, "tourney_hub_guard_")
		}).
		Any()
	assert.True(t, hasHubGuardGroup)
}

func TestWhenTournamentEnabled_SecondPlayerClusterIsUnreachableFromFirst(t *testing.T) {
	t.Parallel()
	perClusterTopologies := []config.MapTopology{
		config.TopologyRing,
		config.TopologyHubAndSpoke,
		config.TopologyCircles,
		config.TopologyChain,
	}
	for _, topology := range perClusterTopologies {
		subTestName := fmt.Sprintf("When%sTopologySelected_SecondPlayerClusterIsUnreachableFromFirst", topology)
		t.Run(subTestName, func(t *testing.T) {
			t.Parallel()
			// Arrange
			generator := test_helpers.NewTemplateGenerator(
				newTournamentConfiguration(topology, 2, gofakeit.Number(1, 20)))

			// Act
			actual, _ := generateTemplate(generator)

			// Assert
			spawnZones := zonesWithPrefix(actual, "Spawn-")
			require.Len(t, spawnZones, 2, "tournament mode must create exactly two spawn zones")

			adjacency := map[string][]string{}
			for _, connection := range actual.Variants[0].Connections {
				adjacency[connection.From] = append(adjacency[connection.From], connection.To)
				adjacency[connection.To] = append(adjacency[connection.To], connection.From)
			}

			reachable := map[string]bool{}
			pending := []string{spawnZones[0].Name}
			for len(pending) > 0 {
				zoneName := pending[len(pending)-1]
				pending = pending[:len(pending)-1]
				if reachable[zoneName] {
					continue
				}
				reachable[zoneName] = true
				pending = append(pending, adjacency[zoneName]...)
			}
			assert.False(t, reachable[spawnZones[1].Name],
				"second player's spawn must not be reachable from the first player's cluster")
		})
	}
}

func TestWhenTournamentEnabledWithRandomPortals_AddsPortalConnections(t *testing.T) {
	t.Parallel()
	// Arrange
	// Portals are drawn from each player's own neutral cluster, so the zone
	// count must stay high enough for both clusters to offer portal targets.
	configuration := newTournamentConfiguration(config.TopologyRing, gofakeit.Number(2, 8), gofakeit.Number(8, 20))
	configuration.RandomPortals = true
	configuration.MaxPortalConnections = 4
	generator := test_helpers.NewTemplateGenerator(configuration)

	// Act
	actual, _ := generateTemplate(generator)

	// Assert
	hasPortalConnections := linq.FromSlice(actual.Variants[0].Connections).
		Where(func(connection entities.Connection) bool { return connection.ConnectionType == "Portal" }).
		Any()
	assert.True(t, hasPortalConnections)
}

// newTournamentConfiguration builds a two-player tournament configuration for
// the given per-cluster topology and neutral zone count.
func newTournamentConfiguration(
	topology config.MapTopology,
	playerCount int,
	neutralZoneCount int) *config.GeneratorConfig {
	configuration := config.NewGeneratorConfig()
	configuration.Topology = topology
	configuration.PlayerCount = playerCount
	configuration.ZoneConfiguration.NeutralZoneCount = neutralZoneCount
	configuration.TournamentRules = &config.TournamentRules{
		Enabled:            true,
		FirstTournamentDay: 14,
		Interval:           7,
		PointsToWin:        2,
	}
	return configuration
}
