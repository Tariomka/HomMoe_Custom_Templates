package templateGenerator_test

import (
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/linq"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenDefaultConfiguration_ReturnsGoldenTemplate(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	generator := test_helpers.NewTemplateGenerator(configuration)
	expected := test_helpers.GetDefaultTemplate()

	// Act
	actual, _ := generateTemplate(generator)

	// Assert
	assert.Equal(t, expected, *actual)
}

func TestWhenTemplateNameIsEmpty_SetsDefaultName(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.TemplateName = ""
	generator := test_helpers.NewTemplateGenerator(configuration)

	// Act
	actual, _ := generateTemplate(generator)

	// Assert
	assert.Equal(t, "Custom Template", actual.Name)
}

func TestWhenTemplateNameProvided_SetsProvidedName(t *testing.T) {
	t.Parallel()
	// Arrange
	expectedName := gofakeit.InputName()
	configuration := config.NewGeneratorConfig()
	configuration.TemplateName = expectedName
	generator := test_helpers.NewTemplateGenerator(configuration)

	// Act
	actual, _ := generateTemplate(generator)

	// Assert
	assert.Equal(t, expectedName, actual.Name)
}

func TestWhenMapSizeProvided_SetsSizeX(t *testing.T) {
	t.Parallel()
	// Arrange
	expectedSize := gofakeit.Number(20, 900)
	configuration := config.NewGeneratorConfig()
	configuration.MapSize = expectedSize
	generator := test_helpers.NewTemplateGenerator(configuration)

	// Act
	actual, _ := generateTemplate(generator)

	// Assert
	assert.Equal(t, expectedSize, actual.SizeX)
}

func TestWhenMapSizeProvided_SetsSizeZ(t *testing.T) {
	t.Parallel()
	// Arrange
	expectedSize := gofakeit.Number(20, 900)
	configuration := config.NewGeneratorConfig()
	configuration.MapSize = expectedSize
	generator := test_helpers.NewTemplateGenerator(configuration)

	// Act
	actual, _ := generateTemplate(generator)

	// Assert
	assert.Equal(t, expectedSize, actual.SizeZ)
}

func TestWhenGameModeProvided_PropagatesGameModeToTemplate(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.GameMode = "SingleHero"
	generator := test_helpers.NewTemplateGenerator(configuration)

	// Act
	actual, _ := generateTemplate(generator)

	// Assert
	assert.Equal(t, "SingleHero", actual.GameMode)
}

func TestWhenVictoryConditionProvided_PropagatesToDisplayWinCondition(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.GameEndConditions = &config.GameEndConditions{
		VictoryCondition: "win_condition_2",
		LostStartCityDay: 3,
		CityHoldDays:     6,
	}
	generator := test_helpers.NewTemplateGenerator(configuration)

	// Act
	actual, _ := generateTemplate(generator)

	// Assert
	assert.Equal(t, "win_condition_2", actual.DisplayWinCondition)
}

func TestWhenGameEndConditionsAreNil_UsesStandardWinCondition(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.GameEndConditions = nil
	generator := test_helpers.NewTemplateGenerator(configuration)

	// Act
	actual, _ := generateTemplate(generator)

	// Assert
	assert.Equal(t, "win_condition_1", actual.DisplayWinCondition)
}

// ── Topology zone structure ──────────────────────────────────────────

func TestWhenRingTopologyAndPlayerCountProvided_CreatesSpawnZonePerPlayer(t *testing.T) {
	t.Parallel()
	// Arrange
	playerCount := gofakeit.Number(2, 8)
	configuration := config.NewGeneratorConfig()
	configuration.PlayerCount = playerCount
	configuration.Topology = config.TopologyRing
	generator := test_helpers.NewTemplateGenerator(configuration)

	// Act
	actual, _ := generateTemplate(generator)

	// Assert
	assert.Len(t, zonesWithPrefix(actual, "Spawn-"), playerCount)
}

func TestWhenRingTopologyAndNeutralZoneCountProvided_CreatesNeutralZonePerCount(t *testing.T) {
	t.Parallel()
	// Arrange
	expectedNeutralZoneCount := gofakeit.Number(0, 30)
	configuration := config.NewGeneratorConfig()
	configuration.ZoneConfiguration.NeutralZoneCount = expectedNeutralZoneCount
	configuration.Topology = config.TopologyRing
	generator := test_helpers.NewTemplateGenerator(configuration)

	// Act
	actual, _ := generateTemplate(generator)

	// Assert
	assert.Len(t, zonesWithPrefix(actual, "Neutral-"), expectedNeutralZoneCount)
}

func TestWhenChainTopologySelected_CreatesZoneCountMinusOneConnections(t *testing.T) {
	t.Parallel()
	// Arrange
	playerCount := gofakeit.Number(2, 8)
	neutralZoneCount := gofakeit.Number(0, 10)
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyChain
	configuration.PlayerCount = playerCount
	configuration.ZoneConfiguration.NeutralZoneCount = neutralZoneCount
	generator := test_helpers.NewTemplateGenerator(configuration)

	// Act
	actual, _ := generateTemplate(generator)

	// Assert
	expectedConnectionCount := playerCount + neutralZoneCount - 1
	assert.Len(t, actual.Variants[0].Connections, expectedConnectionCount)
}

func TestWhenHubAndSpokeTopologySelected_CreatesSingleHubZone(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyHubAndSpoke
	configuration.PlayerCount = gofakeit.Number(2, 8)
	configuration.ZoneConfiguration.NeutralZoneCount = gofakeit.Number(0, 5)
	generator := test_helpers.NewTemplateGenerator(configuration)

	// Act
	actual, _ := generateTemplate(generator)

	// Assert
	hubZones := linq.FromSlice(actual.Variants[0].Zones).
		Where(func(zone entities.Zone) bool { return zone.Name == "Hub" }).
		ToSlice()
	assert.Len(t, hubZones, 1)
}

func TestWhenSharedWebTopologyWithZeroNeutralZones_CreatesOneNeutralZone(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologySharedWeb
	configuration.PlayerCount = gofakeit.Number(3, 8)
	configuration.ZoneConfiguration.NeutralZoneCount = 0
	generator := test_helpers.NewTemplateGenerator(configuration)

	// Act
	actual, _ := generateTemplate(generator)

	// Assert
	assert.Len(t, zonesWithPrefix(actual, "Neutral-"), 1)
}

func TestWhenSharedWebTopologyWithZeroNeutralZones_NamesForcedNeutralZoneAfterPlayers(t *testing.T) {
	t.Parallel()
	// Arrange
	playerCount := gofakeit.Number(3, 8)
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologySharedWeb
	configuration.PlayerCount = playerCount
	configuration.ZoneConfiguration.NeutralZoneCount = 0
	generator := test_helpers.NewTemplateGenerator(configuration)

	// Act
	actual, _ := generateTemplate(generator)

	// Assert
	neutralZones := zonesWithPrefix(actual, "Neutral-")
	expectedName := fmt.Sprintf("Neutral-%c", 'A'+playerCount)
	assert.Equal(t, []string{expectedName}, firstZoneNames(neutralZones))
}

func TestWhenPositionDrivenTopologySelected_SetsGeneratorPositionOnAllZones(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name     string
		topology config.MapTopology
	}{
		{"WhenRandomTopologySelected_SetsGeneratorPositionOnAllZones", config.TopologyRandom},
		{"WhenSquareTopologySelected_SetsGeneratorPositionOnAllZones", config.TopologySquare},
		{"WhenGeometricTopologySelected_SetsGeneratorPositionOnAllZones", config.TopologyGeometric},
		{"WhenCrossTopologySelected_SetsGeneratorPositionOnAllZones", config.TopologyCross},
		{"WhenFractalTopologySelected_SetsGeneratorPositionOnAllZones", config.TopologyFractal},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			// Arrange
			configuration := config.NewGeneratorConfig()
			configuration.Topology = testCase.topology
			configuration.PlayerCount = gofakeit.Number(2, 8)
			configuration.ZoneConfiguration.NeutralZoneCount = gofakeit.Number(0, 5)
			generator := test_helpers.NewTemplateGenerator(configuration)

			// Act
			actual, _ := generateTemplate(generator)

			// Assert
			for _, zone := range actual.Variants[0].Zones {
				assert.NotNil(t, zone.GeneratorPosition, "zone %q should have a generator position", zone.Name)
			}
		})
	}
}

func TestWhenCirclesTopologySelected_SetsGeneratorRingOnAllZones(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyCircles
	configuration.PlayerCount = gofakeit.Number(2, 8)
	configuration.ZoneConfiguration.NeutralZoneCount = gofakeit.Number(0, 5)
	generator := test_helpers.NewTemplateGenerator(configuration)

	// Act
	actual, _ := generateTemplate(generator)

	// Assert
	for _, zone := range actual.Variants[0].Zones {
		assert.NotNil(t, zone.GeneratorRing, "zone %q should have a generator ring", zone.Name)
	}
}

// The fractal layout keeps players apart through its own structure -
// neighbouring fractals meet only at neutral tips - even when
// NoDirectPlayerConnections is off.
func TestWhenFractalTopologySelected_OmitsDirectPlayerConnectionsByDesign(t *testing.T) {
	t.Parallel()
	// Arrange
	playerCount := gofakeit.Number(2, 6)
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyFractal
	configuration.PlayerCount = playerCount
	// Enough neutral zones that every player anchors its own fractal, so the
	// graph stays connected without any direct player-to-player bridge.
	configuration.ZoneConfiguration.NeutralZoneCount = playerCount + gofakeit.Number(3, 8)
	configuration.NoDirectPlayerConnections = false
	generator := test_helpers.NewTemplateGenerator(configuration)

	// Act
	actual, _ := generateTemplate(generator)

	// Assert
	directPlayerConnections := linq.FromSlice(actual.Variants[0].Connections).
		Where(func(connection entities.Connection) bool {
			return connection.ConnectionType == "Direct" &&
				strings.HasPrefix(connection.From, "Spawn-") && strings.HasPrefix(connection.To, "Spawn-")
		}).
		ToSlice()
	assert.Empty(t, directPlayerConnections)
}

// ── Connection-type behaviour ────────────────────────────────────────

func TestWhenRandomPortalsEnabled_AddsPortalConnections(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyRing
	configuration.PlayerCount = 4
	configuration.ZoneConfiguration.NeutralZoneCount = 4
	configuration.RandomPortals = true
	configuration.MaxPortalConnections = gofakeit.Number(1, 8)
	generator := test_helpers.NewTemplateGenerator(configuration)

	// Act
	actual, _ := generateTemplate(generator)

	// Assert
	hasPortalConnections := linq.FromSlice(actual.Variants[0].Connections).
		Where(func(connection entities.Connection) bool { return connection.ConnectionType == "Portal" }).
		Any()
	assert.True(t, hasPortalConnections)
}

func TestWhenRandomPortalsDisabled_AddsNoPortalConnections(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyRing
	configuration.PlayerCount = 4
	configuration.ZoneConfiguration.NeutralZoneCount = 4
	configuration.RandomPortals = false
	generator := test_helpers.NewTemplateGenerator(configuration)

	// Act
	actual, _ := generateTemplate(generator)

	// Assert
	hasPortalConnections := linq.FromSlice(actual.Variants[0].Connections).
		Where(func(connection entities.Connection) bool { return connection.ConnectionType == "Portal" }).
		Any()
	assert.False(t, hasPortalConnections)
}

func TestWhenNoDirectPlayerConnectionsEnabled_OmitsDirectPlayerConnections(t *testing.T) {
	t.Parallel()
	// Arrange
	playerCount := gofakeit.Number(2, 6)
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyRing
	configuration.PlayerCount = playerCount
	configuration.ZoneConfiguration.NeutralZoneCount = playerCount + gofakeit.Number(0, 4)
	configuration.NoDirectPlayerConnections = true
	generator := test_helpers.NewTemplateGenerator(configuration)

	// Act
	actual, _ := generateTemplate(generator)

	// Assert
	// Adjacent players lose their ring edge; connectivity repair may still add
	// guarded Fallback links, so only Ring player-player connections are forbidden.
	directPlayerConnections := linq.FromSlice(actual.Variants[0].Connections).
		Where(func(connection entities.Connection) bool {
			return strings.HasPrefix(connection.Name, "Ring-") &&
				strings.HasPrefix(connection.From, "Spawn-") && strings.HasPrefix(connection.To, "Spawn-")
		}).
		ToSlice()
	assert.Empty(t, directPlayerConnections)
}

// ── Roads ────────────────────────────────────────────────────────────

func TestWhenRoadsEnabled_ProducesRoads(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyRing
	configuration.PlayerCount = gofakeit.Number(2, 8)
	configuration.ZoneConfiguration.NeutralZoneCount = gofakeit.Number(2, 6)
	configuration.GenerateRoads = true
	generator := test_helpers.NewTemplateGenerator(configuration)

	// Act
	actual, _ := generateTemplate(generator)

	// Assert
	hasRoads := linq.FromSlice(actual.Variants[0].Zones).
		Where(func(zone entities.Zone) bool { return len(zone.Roads) > 0 }).
		Any()
	assert.True(t, hasRoads)
}

func TestWhenRoadsDisabled_ProducesNoRoads(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyRing
	configuration.PlayerCount = gofakeit.Number(2, 8)
	configuration.ZoneConfiguration.NeutralZoneCount = gofakeit.Number(2, 6)
	configuration.GenerateRoads = false
	generator := test_helpers.NewTemplateGenerator(configuration)

	// Act
	actual, _ := generateTemplate(generator)

	// Assert
	zonesWithRoads := linq.FromSlice(actual.Variants[0].Zones).
		Where(func(zone entities.Zone) bool { return len(zone.Roads) > 0 }).
		ToSlice()
	assert.Empty(t, zonesWithRoads)
}

// ── Castle factions ──────────────────────────────────────────────────

func TestWhenMatchPlayerCastleFactionsEnabled_SetsMatchFactionOnExtraPlayerCastles(t *testing.T) {
	t.Parallel()
	// Arrange
	playerCount := gofakeit.Number(2, 8)
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyRing
	configuration.PlayerCount = playerCount
	configuration.ZoneConfiguration.PlayerZoneCastles = 2
	configuration.MatchPlayerCastleFactions = true
	generator := test_helpers.NewTemplateGenerator(configuration)

	// Act
	actual, _ := generateTemplate(generator)

	// Assert
	expectedFactionTypes := slices.Repeat([]string{"Match"}, playerCount)
	assert.Equal(t, expectedFactionTypes, extraCastleFactionTypes(zonesWithPrefix(actual, "Spawn-")))
}

func TestWhenMatchPlayerCastleFactionsDisabled_SetsRandomFactionOnExtraPlayerCastles(t *testing.T) {
	t.Parallel()
	// Arrange
	playerCount := gofakeit.Number(2, 8)
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyRing
	configuration.PlayerCount = playerCount
	configuration.ZoneConfiguration.PlayerZoneCastles = 2
	configuration.MatchPlayerCastleFactions = false
	generator := test_helpers.NewTemplateGenerator(configuration)

	// Act
	actual, _ := generateTemplate(generator)

	// Assert
	expectedFactionTypes := slices.Repeat([]string{"Random"}, playerCount)
	assert.Equal(t, expectedFactionTypes, extraCastleFactionTypes(zonesWithPrefix(actual, "Spawn-")))
}

// ── City hold / hold-city objects ────────────────────────────────────

func TestWhenCityHoldEnabled_MarksHoldCityWinConditionObjectInZones(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyHubAndSpoke
	configuration.PlayerCount = gofakeit.Number(2, 8)
	configuration.ZoneConfiguration.NeutralZoneCount = gofakeit.Number(1, 6)
	configuration.GameEndConditions = &config.GameEndConditions{
		VictoryCondition: "win_condition_1",
		CityHold:         true,
		CityHoldDays:     gofakeit.Number(1, 10),
		LostStartCityDay: 3,
	}
	generator := test_helpers.NewTemplateGenerator(configuration)

	// Act
	actual, _ := generateTemplate(generator)

	// Assert
	holdCityZones := linq.FromSlice(actual.Variants[0].Zones).
		Where(func(zone entities.Zone) bool {
			return linq.FromSlice(zone.MainObjects).
				Where(func(mainObject entities.MainObject) bool { return mainObject.HoldCityWinCon }).
				Any()
		}).
		ToSlice()
	assert.NotEmpty(t, holdCityZones)
}

func TestWhenCityHoldEnabledWithHubAndSpokeTopology_MarksHubAsHoldCity(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyHubAndSpoke
	configuration.PlayerCount = gofakeit.Number(2, 8)
	configuration.GameEndConditions = &config.GameEndConditions{
		VictoryCondition: "win_condition_5",
		LostStartCityDay: 3,
		CityHoldDays:     6,
		CityHold:         true,
	}
	generator := test_helpers.NewTemplateGenerator(configuration)

	// Act
	actual, _ := generateTemplate(generator)

	// Assert
	hubHoldsCity := linq.FromSlice(actual.Variants[0].Zones).
		Where(func(zone entities.Zone) bool { return zone.Name == "Hub" }).
		Where(func(zone entities.Zone) bool {
			return linq.FromSlice(zone.MainObjects).
				Where(func(mainObject entities.MainObject) bool { return mainObject.HoldCityWinCon }).
				Any()
		}).
		Any()
	assert.True(t, hubHoldsCity)
}

// ── Tournament topology variants ─────────────────────────────────────

func TestWhenTournamentEnabledWithTwoPlayersAndRingTopology_CreatesRingGuardGroups(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyRing
	configuration.PlayerCount = 2
	configuration.ZoneConfiguration.NeutralZoneCount = gofakeit.Number(2, 6)
	configuration.TournamentRules = &config.TournamentRules{
		Enabled:            true,
		FirstTournamentDay: 14,
		Interval:           7,
		PointsToWin:        2,
	}
	generator := test_helpers.NewTemplateGenerator(configuration)

	// Act
	actual, _ := generateTemplate(generator)

	// Assert
	hasRingGuardGroup := linq.FromSlice(actual.Variants[0].Connections).
		Where(func(connection entities.Connection) bool {
			return strings.HasPrefix(connection.GuardMatchGroup, "tourney_ring_guard_")
		}).
		Any()
	assert.True(t, hasRingGuardGroup)
}

func TestWhenTournamentEnabledWithHubAndSpokeTopology_CreatesHubPerPlayer(t *testing.T) {
	t.Parallel()
	// Arrange
	const expectedHubCount = 2 // Tournament mode is only triggered for exactly 2 players.
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyHubAndSpoke
	configuration.PlayerCount = expectedHubCount
	configuration.ZoneConfiguration.NeutralZoneCount = gofakeit.Number(2, 6)
	configuration.TournamentRules = &config.TournamentRules{
		Enabled:            true,
		FirstTournamentDay: 14,
		Interval:           7,
		PointsToWin:        2,
	}
	generator := test_helpers.NewTemplateGenerator(configuration)

	// Act
	actual, _ := generateTemplate(generator)

	// Assert
	assert.Len(t, zonesWithPrefix(actual, "Hub-"), expectedHubCount)
}

func TestWhenTournamentEnabledWithChainTopology_CreatesChainGuardGroups(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyChain
	configuration.PlayerCount = 2
	configuration.ZoneConfiguration.NeutralZoneCount = gofakeit.Number(2, 6)
	configuration.TournamentRules = &config.TournamentRules{
		Enabled:            true,
		FirstTournamentDay: 14,
		Interval:           7,
		PointsToWin:        2,
	}
	generator := test_helpers.NewTemplateGenerator(configuration)

	// Act
	actual, _ := generateTemplate(generator)

	// Assert
	hasChainGuardGroup := linq.FromSlice(actual.Variants[0].Connections).
		Where(func(connection entities.Connection) bool {
			return strings.HasPrefix(connection.GuardMatchGroup, "tourney_guard_")
		}).
		Any()
	assert.True(t, hasChainGuardGroup)
}

func TestWhenTournamentEnabledWithCirclesTopology_CreatesBalancedGuardGroups(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyCircles
	configuration.PlayerCount = 2
	configuration.ZoneConfiguration.NeutralZoneCount = gofakeit.Number(2, 6)
	configuration.TournamentRules = &config.TournamentRules{
		Enabled:            true,
		FirstTournamentDay: 14,
		Interval:           7,
		PointsToWin:        2,
	}
	generator := test_helpers.NewTemplateGenerator(configuration)

	// Act
	actual, _ := generateTemplate(generator)

	// Assert
	hasBalancedGuardGroup := linq.FromSlice(actual.Variants[0].Connections).
		Where(func(connection entities.Connection) bool {
			return strings.HasPrefix(connection.GuardMatchGroup, "tourney_bal_guard_")
		}).
		Any()
	assert.True(t, hasBalancedGuardGroup)
}

// ── Advanced neutral mix ─────────────────────────────────────────────

func TestWhenAdvancedNeutralMixEnabled_CreatesNeutralZonePerConfiguredTierCount(t *testing.T) {
	t.Parallel()
	// Arrange
	remainingBudget := 30
	lowNoCastleCount := gofakeit.Number(0, remainingBudget)
	remainingBudget -= lowNoCastleCount
	lowCastleCount := gofakeit.Number(0, remainingBudget)
	remainingBudget -= lowCastleCount
	mediumNoCastleCount := gofakeit.Number(0, remainingBudget)
	remainingBudget -= mediumNoCastleCount
	mediumCastleCount := gofakeit.Number(0, remainingBudget)
	remainingBudget -= mediumCastleCount
	highNoCastleCount := gofakeit.Number(0, remainingBudget)
	remainingBudget -= highNoCastleCount
	highCastleCount := gofakeit.Number(0, remainingBudget)
	expectedNeutralCount := lowNoCastleCount + lowCastleCount + mediumNoCastleCount +
		mediumCastleCount + highNoCastleCount + highCastleCount

	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyRing
	configuration.ZoneConfiguration.Advanced.Enabled = true
	configuration.ZoneConfiguration.Advanced.NeutralLowNoCastleCount = lowNoCastleCount
	configuration.ZoneConfiguration.Advanced.NeutralLowCastleCount = lowCastleCount
	configuration.ZoneConfiguration.Advanced.NeutralMediumNoCastleCount = mediumNoCastleCount
	configuration.ZoneConfiguration.Advanced.NeutralMediumCastleCount = mediumCastleCount
	configuration.ZoneConfiguration.Advanced.NeutralHighNoCastleCount = highNoCastleCount
	configuration.ZoneConfiguration.Advanced.NeutralHighCastleCount = highCastleCount
	generator := test_helpers.NewTemplateGenerator(configuration)

	// Act
	actual, _ := generateTemplate(generator)

	// Assert
	assert.Len(t, zonesWithPrefix(actual, "Neutral-"), expectedNeutralCount)
}

func TestWhenAdvancedGuardRandomizationExceedsMaximum_ClampsGuardRandomization(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.ZoneConfiguration.Advanced.Enabled = true
	configuration.ZoneConfiguration.GuardRandomization = gofakeit.Float64Range(0.6, 10.0)
	generator := test_helpers.NewTemplateGenerator(configuration)

	// Act
	actual, _ := generateTemplate(generator)

	// Assert
	for _, zone := range actual.Variants[0].Zones {
		assert.LessOrEqual(t, zone.GuardRandomization, 0.5, "zone %q guard randomization should be clamped", zone.Name)
	}
}

// ── Structural template fields ───────────────────────────────────────

func TestWhenGenerating_ProducesExactlyOneVariant(t *testing.T) {
	t.Parallel()
	// Arrange
	generator := test_helpers.NewTemplateGenerator(config.NewGeneratorConfig())

	// Act
	actual, _ := generateTemplate(generator)

	// Assert
	assert.Len(t, actual.Variants, 1)
}

func TestWhenGenerating_ProducesFourZoneLayouts(t *testing.T) {
	t.Parallel()
	// Arrange
	generator := test_helpers.NewTemplateGenerator(config.NewGeneratorConfig())

	// Act
	actual, _ := generateTemplate(generator)

	// Assert
	assert.Len(t, actual.ZoneLayouts, 4)
}

func TestWhenGenerating_ProducesContentCountLimits(t *testing.T) {
	t.Parallel()
	// Arrange
	generator := test_helpers.NewTemplateGenerator(config.NewGeneratorConfig())

	// Act
	actual, _ := generateTemplate(generator)

	// Assert
	assert.NotEmpty(t, actual.ContentCountLimits)
}

// ── Description ──────────────────────────────────────────────────────

func TestWhenChainTopologySelected_IncludesTopologyNameInDescription(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyChain
	configuration.PlayerCount = gofakeit.Number(2, 8)
	configuration.ZoneConfiguration.NeutralZoneCount = gofakeit.Number(0, 5)
	generator := test_helpers.NewTemplateGenerator(configuration)

	// Act
	actual, _ := generateTemplate(generator)

	// Assert
	assert.Contains(t, actual.Description, "Chain")
}

func TestWhenDescriptionOptionsEnabled_AppendsOptionPhrases(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyRing
	configuration.PlayerCount = gofakeit.Number(2, 8)
	configuration.ZoneConfiguration.NeutralZoneCount = gofakeit.Number(2, 6)
	configuration.NoDirectPlayerConnections = true
	configuration.RandomPortals = true
	configuration.SpawnRemoteFootholds = false
	configuration.GenerateRoads = false
	generator := test_helpers.NewTemplateGenerator(configuration)

	// Act
	actual, _ := generateTemplate(generator)

	// Assert
	cases := []struct {
		name   string
		phrase string
	}{
		{"WhenNoDirectPlayerConnectionsEnabled_MentionsIsolatedPlayerStarts", "isolated player starts"},
		{"WhenRandomPortalsEnabled_MentionsRandomPortals", "random portals"},
		{"WhenRemoteFootholdsDisabled_MentionsNoRemoteFootholds", "no remote footholds"},
		{"WhenRoadsDisabled_MentionsRoadsDisabled", "roads disabled"},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			assert.Contains(t, actual.Description, testCase.phrase)
		})
	}
}

// ── Mandatory content groups ─────────────────────────────────────────

func TestWhenGenerating_CreatesMandatoryContentGroupPerPlayer(t *testing.T) {
	t.Parallel()
	// Arrange
	playerCount := gofakeit.Number(2, 8)
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyRing
	configuration.PlayerCount = playerCount
	configuration.ZoneConfiguration.NeutralZoneCount = gofakeit.Number(1, 6)
	generator := test_helpers.NewTemplateGenerator(configuration)

	// Act
	actual, _ := generateTemplate(generator)

	// Assert
	playerGroups := linq.FromSlice(actual.MandatoryContent).
		Where(func(group entities.MandatoryContent) bool {
			return strings.HasPrefix(group.Name, "mandatory_content_side_")
		}).
		ToSlice()
	assert.Len(t, playerGroups, playerCount)
}

func TestWhenGenerating_CreatesMandatoryContentGroupPerNeutralZone(t *testing.T) {
	t.Parallel()
	// Arrange
	neutralZoneCount := gofakeit.Number(1, 6)
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyRing
	configuration.PlayerCount = gofakeit.Number(2, 8)
	configuration.ZoneConfiguration.NeutralZoneCount = neutralZoneCount
	generator := test_helpers.NewTemplateGenerator(configuration)

	// Act
	actual, _ := generateTemplate(generator)

	// Assert
	neutralGroups := linq.FromSlice(actual.MandatoryContent).
		Where(func(group entities.MandatoryContent) bool {
			return strings.HasPrefix(group.Name, "mandatory_content_neutral_")
		}).
		ToSlice()
	assert.Len(t, neutralGroups, neutralZoneCount)
}

// ── Spawn zone main objects ──────────────────────────────────────────

func TestWhenGenerating_PlacesSpawnMainObjectFirstInEachSpawnZone(t *testing.T) {
	t.Parallel()
	// Arrange
	playerCount := gofakeit.Number(2, 8)
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyRing
	configuration.PlayerCount = playerCount
	generator := test_helpers.NewTemplateGenerator(configuration)

	// Act
	actual, _ := generateTemplate(generator)

	// Assert
	expectedTypes := slices.Repeat([]string{"Spawn"}, playerCount)
	assert.Equal(t, expectedTypes, firstMainObjectTypes(zonesWithPrefix(actual, "Spawn-")))
}

// zonesWithPrefix returns the zones of the template's first variant whose name
// starts with the given prefix.
func zonesWithPrefix(generated *entities.RmgTemplate, prefix string) []entities.Zone {
	return linq.FromSlice(generated.Variants[0].Zones).
		Where(func(zone entities.Zone) bool { return strings.HasPrefix(zone.Name, prefix) }).
		ToSlice()
}

// extraCastleFactionTypes collects the faction type of the second main object
// (the first extra castle) of every spawn zone, using "<missing>" when the
// castle or its faction is absent so a mismatch shows up in the assertion.
func extraCastleFactionTypes(spawnZones []entities.Zone) []string {
	var factionTypes []string
	for _, zone := range spawnZones {
		if len(zone.MainObjects) < 2 || zone.MainObjects[1].Faction == nil {
			factionTypes = append(factionTypes, "<missing>")
			continue
		}
		factionTypes = append(factionTypes, zone.MainObjects[1].Faction.Type)
	}
	return factionTypes
}

// firstMainObjectTypes collects the type of the first main object of every
// given zone, using "<missing>" when a zone has no main objects.
func firstMainObjectTypes(zones []entities.Zone) []string {
	var objectTypes []string
	for _, zone := range zones {
		if len(zone.MainObjects) == 0 {
			objectTypes = append(objectTypes, "<missing>")
			continue
		}
		objectTypes = append(objectTypes, zone.MainObjects[0].Type)
	}
	return objectTypes
}

// firstZoneNames returns the names of the given zones in order.
func firstZoneNames(zones []entities.Zone) []string {
	var names []string
	for _, zone := range zones {
		names = append(names, zone.Name)
	}
	return names
}
