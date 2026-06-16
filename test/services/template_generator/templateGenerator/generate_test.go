package templateGenerator_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/linq"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenValidConfig_ReturnsExpectedResult(t *testing.T) {
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.ShufflePlayerZones = false // Deterministic player-zone ordering for a stable golden comparison.
	templateGenerator := template_generator.NewTemplateGenerator(configuration)
	expected := test_helpers.GetDefaultTemplate()

	// Act
	actual := templateGenerator.Generate()

	// Assert
	assert.Equal(t, expected, *actual)
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, expected.GameMode, actual.GameMode)
	assert.Equal(t, expected.Description, actual.Description)
	assert.Equal(t, expected.SizeX, actual.SizeX)
	assert.Equal(t, expected.SizeZ, actual.SizeZ)
	assert.Equal(t, expected.ValueOverrides, actual.ValueOverrides)
	assert.Equal(t, expected.Orientation, actual.Orientation)
	assert.Equal(t, expected.Border, actual.Border)
	assert.Equal(t, expected.GameRules, actual.GameRules)
	assert.Equal(t, expected.GlobalBans, actual.GlobalBans)
	assert.Equal(t, expected.Variants[0].Orientation, actual.Variants[0].Orientation)
	assert.Equal(t, expected.Variants[0].Connections, actual.Variants[0].Connections)
	assert.Equal(t, expected.Variants[0].Border, actual.Variants[0].Border)
	assert.Equal(t, expected.Variants[0].Zones, actual.Variants[0].Zones)
	assert.Equal(t, expected.ZoneLayouts, actual.ZoneLayouts)
	assert.Equal(t, expected.MandatoryContent, actual.MandatoryContent)
	assert.Equal(t, expected.ContentCountLimits, actual.ContentCountLimits)
	assert.Equal(t, expected.ContentPools, actual.ContentPools)
	assert.Equal(t, expected.ContentLists, actual.ContentLists)
}

func TestWhenTemplateNameIsEmpty_SetsDefaultName(t *testing.T) {
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.TemplateName = ""
	templateGenerator := template_generator.NewTemplateGenerator(configuration)

	// Act
	actual := templateGenerator.Generate()

	// Assert
	assert.Equal(t, actual.Name, "Custom Template")
}

func TestWhenNonDefaultTemplateName_SetsExpectedName(t *testing.T) {
	// Arrange
	expectedName := gofakeit.InputName()
	configuration := config.NewGeneratorConfig()
	configuration.TemplateName = expectedName
	templateGenerator := template_generator.NewTemplateGenerator(configuration)

	// Act
	actual := templateGenerator.Generate()

	// Assert
	assert.Equal(t, actual.Name, expectedName)
}

func TestWhenNonDefaultMapSize_SetsExpectedSize(t *testing.T) {
	// Arrange
	expectedSize := gofakeit.Number(20, 900)
	configuration := config.NewGeneratorConfig()
	configuration.MapSize = expectedSize
	templateGenerator := template_generator.NewTemplateGenerator(configuration)

	// Act
	actual := templateGenerator.Generate()

	// Assert
	assert.Equal(t, expectedSize, actual.SizeX)
	assert.Equal(t, expectedSize, actual.SizeZ)
}

func TestWhenGameModeIsClassic_SetsExpectedGameModeAndDisablesHeroHireBan(t *testing.T) {
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.GameMode = "Classic"
	templateGenerator := template_generator.NewTemplateGenerator(configuration)

	// Act
	actual := templateGenerator.Generate()

	// Assert
	assert.Equal(t, "Classic", actual.GameMode)
	assert.False(t, actual.GameRules.HeroHireBan)
}

func TestWhenGameModeIsSingleHero_SetsExpectedHeroCounts(t *testing.T) {
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.GameMode = "SingleHero"
	templateGenerator := template_generator.NewTemplateGenerator(configuration)

	// Act
	actual := templateGenerator.Generate()

	// Assert
	assert.Equal(t, "SingleHero", actual.GameMode)
	assert.True(t, actual.GameRules.HeroHireBan)
	assert.Equal(t, 1, actual.GameRules.HeroCountMin)
	assert.Equal(t, 1, actual.GameRules.HeroCountMax)
	assert.Equal(t, 1, actual.GameRules.HeroCountIncrement)
}

func TestWhenNonDefaultHeroSettingsInClassic_SetsExpectedHeroSettings(t *testing.T) {
	// Arrange
	expectedMin := gofakeit.Number(1, 5)
	expectedMax := gofakeit.Number(expectedMin, 10)
	expectedIncrement := gofakeit.Number(1, 3)
	configuration := config.NewGeneratorConfig()
	configuration.HeroSettings = config.HeroSettings{
		HeroCountMin:       expectedMin,
		HeroCountMax:       expectedMax,
		HeroCountIncrement: expectedIncrement,
	}
	templateGenerator := template_generator.NewTemplateGenerator(configuration)

	// Act
	actual := templateGenerator.Generate()

	// Assert
	assert.Equal(t, expectedIncrement, actual.GameRules.HeroCountIncrement)
	assert.Equal(t, expectedMax, actual.GameRules.HeroCountMax)
	assert.Equal(t, expectedMin, actual.GameRules.HeroCountMin)
}

func TestWhenNonDefaultFactionLawsExpPercent_SetsExpectedModifier(t *testing.T) {
	// Arrange
	expected := gofakeit.IntRange(25, 200)
	configuration := config.NewGeneratorConfig()
	configuration.FactionLawsExpPercent = expected
	templateGenerator := template_generator.NewTemplateGenerator(configuration)

	// Act
	actual := templateGenerator.Generate()

	// Assert
	assert.Equal(t, float64(expected)/100, actual.GameRules.FactionLawsExpModifier)
}

func TestWhenBelowMinimumFactionLawsExpPercent_SetsExpectedModifier(t *testing.T) {
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.FactionLawsExpPercent = gofakeit.IntRange(-1000, 20)
	templateGenerator := template_generator.NewTemplateGenerator(configuration)

	// Act
	actual := templateGenerator.Generate()

	// Assert
	assert.Equal(t, 0.25, actual.GameRules.FactionLawsExpModifier)
}

func TestWhenAboveMaximumFactionLawsExpPercent_SetsExpectedModifier(t *testing.T) {
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.FactionLawsExpPercent = gofakeit.IntRange(201, 9999)
	templateGenerator := template_generator.NewTemplateGenerator(configuration)

	// Act
	actual := templateGenerator.Generate()

	// Assert
	assert.Equal(t, 2.0, actual.GameRules.FactionLawsExpModifier)
}

func TestWhenNonDefaultAstrologyExpPercent_SetsExpectedModifier(t *testing.T) {
	// Arrange
	expected := gofakeit.IntRange(25, 200)
	configuration := config.NewGeneratorConfig()
	configuration.AstrologyExpPercent = expected
	templateGenerator := template_generator.NewTemplateGenerator(configuration)

	// Act
	actual := templateGenerator.Generate()

	// Assert
	assert.Equal(t, float64(expected)/100, actual.GameRules.AstrologyExpModifier)
}

func TestWhenBelowMinimumAstrologyExpPercent_SetsExpectedModifier(t *testing.T) {
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.AstrologyExpPercent = gofakeit.IntRange(-1000, 20)
	templateGenerator := template_generator.NewTemplateGenerator(configuration)

	// Act
	actual := templateGenerator.Generate()

	// Assert
	assert.Equal(t, 0.25, actual.GameRules.AstrologyExpModifier)
}

func TestWhenAboveMaximumAstrologyExpPercent_SetsExpectedModifier(t *testing.T) {
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.AstrologyExpPercent = gofakeit.IntRange(201, 9999)
	templateGenerator := template_generator.NewTemplateGenerator(configuration)

	// Act
	actual := templateGenerator.Generate()

	// Assert
	assert.Equal(t, 2.0, actual.GameRules.AstrologyExpModifier)
}

func TestWhenDefaultTopologySelectedAndNonDefaultPlayerCountProvided_SetsExpectedSpawnZones(t *testing.T) {
	// Arrange
	playerCount := gofakeit.Number(2, 8)
	configuration := config.NewGeneratorConfig()
	configuration.PlayerCount = playerCount
	configuration.Topology = config.TopologyDefault
	templateGenerator := template_generator.NewTemplateGenerator(configuration)

	// Act
	actual := templateGenerator.Generate()

	// Assert
	spawnZones := linq.FromSlice(actual.Variants[0].Zones).
		Where(func(x entities.Zone) bool { return strings.HasPrefix(x.Name, "Spawn-") }).
		ToSlice()
	assert.Equal(t, playerCount, len(spawnZones))
	// for i, zone := range spawnZones {
	// 	expectedName := fmt.Sprintf("Spawn-%c", 'A'+i)
	// 	assert.Equal(t, expectedName, zone.Name)
	// }
}
func TestWhenDefaultTopologySelectedAndNonDefaultNeutralZoneCountProvided_SetsExpectedNeutralZones(t *testing.T) {
	// Arrange
	expectedNeutralZoneCount := gofakeit.Number(0, 30)
	configuration := config.NewGeneratorConfig()
	configuration.ZoneConfiguration.NeutralZoneCount = expectedNeutralZoneCount
	configuration.Topology = config.TopologyDefault
	templateGenerator := template_generator.NewTemplateGenerator(configuration)

	// Act
	actual := templateGenerator.Generate()

	// Assert
	neutralZones := linq.FromSlice(actual.Variants[0].Zones).
		Where(func(x entities.Zone) bool { return strings.HasPrefix(x.Name, "Neutral-") }).
		ToSlice()
	assert.Equal(t, expectedNeutralZoneCount, len(neutralZones))
	// for i, zone := range neutralZones {
	// 	expectedName := fmt.Sprintf("Neutral-%c", 'A'+i+2) // default 2 player zones
	// 	assert.Equal(t, expectedName, zone.Name)
	// }
}

func TestWhenChainTopologySelected_HasZoneCountMinusOneConnections(t *testing.T) {
	// Arrange
	playerCount := gofakeit.Number(2, 8)
	neutralZoneCount := gofakeit.Number(0, 10)
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyChain
	configuration.PlayerCount = playerCount
	configuration.ZoneConfiguration.NeutralZoneCount = neutralZoneCount
	templateGenerator := template_generator.NewTemplateGenerator(configuration)

	// Act
	actual := templateGenerator.Generate()

	// Assert
	expectedConnectionCount := playerCount + neutralZoneCount - 1
	assert.Equal(t, expectedConnectionCount, len(actual.Variants[0].Connections))
}

func TestWhenHubTopologySelected_CreatesHubZone(t *testing.T) {
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyHubAndSpoke
	configuration.PlayerCount = gofakeit.Number(2, 8)
	configuration.ZoneConfiguration.NeutralZoneCount = gofakeit.Number(0, 5)
	templateGenerator := template_generator.NewTemplateGenerator(configuration)

	// Act
	actual := templateGenerator.Generate()

	// Assert
	hubZones := linq.FromSlice(actual.Variants[0].Zones).
		Where(func(x entities.Zone) bool { return x.Name == "Hub" }).
		ToSlice()
	assert.Equal(t, 1, len(hubZones))
}

func TestWhenSharedWebTopologySelectedAndZeroNeutralZonesProvided_SetsMinimumNeutralZones(t *testing.T) {
	// Arrange
	playerCount := gofakeit.Number(3, 8)
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologySharedWeb
	configuration.PlayerCount = playerCount
	configuration.ZoneConfiguration.NeutralZoneCount = 0
	templateGenerator := template_generator.NewTemplateGenerator(configuration)

	// Act
	actual := templateGenerator.Generate()

	// Assert
	actualNeutralZones := linq.FromSlice(actual.Variants[0].Zones).
		Where(func(x entities.Zone) bool { return strings.HasPrefix(x.Name, "Neutral-") }).
		ToSlice()
	assert.Equal(t, 1, len(actualNeutralZones))
	assert.Equal(t, fmt.Sprintf("Neutral-%c", 'A'+playerCount), actualNeutralZones[0].Name)
}

func TestWhenRandomTopologySelected_SetsGeneratorPositionOnAllZones(t *testing.T) {
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyRandom
	configuration.PlayerCount = gofakeit.Number(2, 8)
	configuration.ZoneConfiguration.NeutralZoneCount = gofakeit.Number(0, 5)
	templateGenerator := template_generator.NewTemplateGenerator(configuration)

	// Act
	actual := templateGenerator.Generate()

	// Assert
	for _, zone := range actual.Variants[0].Zones {
		assert.NotNil(t, zone.GeneratorPosition, "zone %q should have a generator position", zone.Name)
	}
}

func TestWhenCirclesTopologySelected_SetsGeneratorRingOnAllZones(t *testing.T) {
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyCircles
	configuration.PlayerCount = gofakeit.Number(2, 8)
	configuration.ZoneConfiguration.NeutralZoneCount = gofakeit.Number(0, 5)
	templateGenerator := template_generator.NewTemplateGenerator(configuration)

	// Act
	actual := templateGenerator.Generate()

	// Assert
	for _, zone := range actual.Variants[0].Zones {
		assert.NotNil(t, zone.GeneratorRing, "zone %q should have a generator ring", zone.Name)
	}
}

func TestWhenSquareTopologySelected_SetsGeneratorPositionOnAllZones(t *testing.T) {
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologySquare
	configuration.PlayerCount = gofakeit.Number(2, 8)
	configuration.ZoneConfiguration.NeutralZoneCount = gofakeit.Number(0, 5)
	templateGenerator := template_generator.NewTemplateGenerator(configuration)

	// Act
	actual := templateGenerator.Generate()

	// Assert
	for _, zone := range actual.Variants[0].Zones {
		assert.NotNil(t, zone.GeneratorPosition, "zone %q should have a generator position", zone.Name)
	}
}

func TestWhenGeometricTopologySelected_SetsGeneratorPositionOnAllZones(t *testing.T) {
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyGeometric
	configuration.PlayerCount = gofakeit.Number(2, 8)
	configuration.ZoneConfiguration.NeutralZoneCount = gofakeit.Number(0, 5)
	templateGenerator := template_generator.NewTemplateGenerator(configuration)

	// Act
	actual := templateGenerator.Generate()

	// Assert
	for _, zone := range actual.Variants[0].Zones {
		assert.NotNil(t, zone.GeneratorPosition, "zone %q should have a generator position", zone.Name)
	}
}

func TestWhenCrossTopologySelected_SetsGeneratorPositionOnAllZones(t *testing.T) {
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyCross
	configuration.PlayerCount = gofakeit.Number(2, 8)
	configuration.ZoneConfiguration.NeutralZoneCount = gofakeit.Number(0, 5)
	templateGenerator := template_generator.NewTemplateGenerator(configuration)

	// Act
	actual := templateGenerator.Generate()

	// Assert
	for _, zone := range actual.Variants[0].Zones {
		assert.NotNil(t, zone.GeneratorPosition, "zone %q should have a generator position", zone.Name)
	}
}

func TestWhenFractalTopologySelected_SetsGeneratorPositionOnAllZones(t *testing.T) {
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyFractal
	configuration.PlayerCount = gofakeit.Number(2, 8)
	configuration.ZoneConfiguration.NeutralZoneCount = gofakeit.Number(0, 5)
	templateGenerator := template_generator.NewTemplateGenerator(configuration)

	// Act
	actual := templateGenerator.Generate()

	// Assert
	for _, zone := range actual.Variants[0].Zones {
		assert.NotNil(t, zone.GeneratorPosition, "zone %q should have a generator position", zone.Name)
	}
}

// TestWhenFractalTopologySelected_OmitsDirectPlayerConnectionsByDesign verifies
// the fractal layout keeps players apart through its own structure — neighbouring
// fractals meet only at neutral tips — even when NoDirectPlayerConnections is off.
func TestWhenFractalTopologySelected_OmitsDirectPlayerConnectionsByDesign(t *testing.T) {
	// Arrange
	playerCount := gofakeit.Number(2, 6)
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyFractal
	configuration.PlayerCount = playerCount
	// Enough neutral zones that every player anchors its own fractal, so the
	// graph stays connected without any direct player-to-player bridge.
	configuration.ZoneConfiguration.NeutralZoneCount = playerCount + gofakeit.Number(3, 8)
	configuration.NoDirectPlayerConnections = false
	templateGenerator := template_generator.NewTemplateGenerator(configuration)

	// Act
	actual := templateGenerator.Generate()

	// Assert
	directPlayerConnections := linq.FromSlice(actual.Variants[0].Connections).
		Where(func(x entities.Connection) bool {
			return x.ConnectionType == "Direct" &&
				strings.HasPrefix(x.From, "Spawn-") && strings.HasPrefix(x.To, "Spawn-")
		}).
		ToSlice()
	assert.Empty(t, directPlayerConnections)
}

// ── Generate: connection-type behaviour ──────────────────────────────

func TestWhenRandomPortalsEnabled_AddsPortalConnections(t *testing.T) {
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyDefault
	configuration.PlayerCount = 4
	configuration.ZoneConfiguration.NeutralZoneCount = 4
	configuration.RandomPortals = true
	configuration.MaxPortalConnections = gofakeit.Number(1, 8)
	templateGenerator := template_generator.NewTemplateGenerator(configuration)

	// Act
	actual := templateGenerator.Generate()

	// Assert
	hasPortalConnections := linq.FromSlice(actual.Variants[0].Connections).
		Where(func(x entities.Connection) bool { return x.ConnectionType == "Portal" }).
		Any()
	assert.True(t, hasPortalConnections)
}

func TestWhenRandomPortalsDisabled_AddsNoPortalConnections(t *testing.T) {
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyDefault
	configuration.PlayerCount = 4
	configuration.ZoneConfiguration.NeutralZoneCount = 4
	configuration.RandomPortals = false
	templateGenerator := template_generator.NewTemplateGenerator(configuration)

	// Act
	actual := templateGenerator.Generate()

	// Assert
	hasPortalConnections := linq.FromSlice(actual.Variants[0].Connections).
		Where(func(x entities.Connection) bool { return x.ConnectionType == "Portal" }).
		Any()
	assert.False(t, hasPortalConnections)
}

func TestWhenNoDirectPlayerConnectionsEnabled_OmitsDirectPlayerConnections(t *testing.T) {
	// Arrange
	playerCount := gofakeit.Number(2, 6)
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyDefault
	configuration.PlayerCount = playerCount
	// Provide at least one neutral zone between every player so the topology can
	// actually honour the separation instead of falling back to direct links.
	configuration.MinNeutralZonesBetweenPlayers = 1
	configuration.ZoneConfiguration.NeutralZoneCount = playerCount + gofakeit.Number(0, 4)

	configuration.NoDirectPlayerConnections = true
	templateGenerator := template_generator.NewTemplateGenerator(configuration)

	// Act
	actual := templateGenerator.Generate()

	// Assert
	directPlayerConnections := linq.FromSlice(actual.Variants[0].Connections).
		Where(func(x entities.Connection) bool {
			return x.ConnectionType == "Direct" &&
				strings.HasPrefix(x.From, "Spawn-") && strings.HasPrefix(x.To, "Spawn-")
		}).
		ToSlice()
	assert.Empty(t, directPlayerConnections)
}

// ── Generate: roads ──────────────────────────────────────────────────

func TestWhenRoadsEnabled_ProducesRoads(t *testing.T) {
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyDefault
	configuration.PlayerCount = gofakeit.Number(2, 8)
	configuration.ZoneConfiguration.NeutralZoneCount = gofakeit.Number(2, 6)
	configuration.GenerateRoads = true
	templateGenerator := template_generator.NewTemplateGenerator(configuration)

	// Act
	actual := templateGenerator.Generate()

	// Assert
	hasRoads := linq.FromSlice(actual.Variants[0].Zones).
		Where(func(x entities.Zone) bool { return len(x.Roads) > 0 }).
		Any()
	assert.True(t, hasRoads)
}

func TestWhenRoadsDisabled_ProducesNoRoads(t *testing.T) {
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyDefault
	configuration.PlayerCount = gofakeit.Number(2, 8)
	configuration.ZoneConfiguration.NeutralZoneCount = gofakeit.Number(2, 6)
	configuration.GenerateRoads = false
	templateGenerator := template_generator.NewTemplateGenerator(configuration)

	// Act
	actual := templateGenerator.Generate()

	// Assert
	for _, zone := range actual.Variants[0].Zones {
		assert.Empty(t, zone.Roads, "zone %q should have no roads", zone.Name)
	}
}

// ── Generate: castle factions ────────────────────────────────────────

func TestWhenMatchPlayerCastleFactionsEnabled_SetsMatchFactionOnPlayerCastles(t *testing.T) {
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyDefault
	configuration.PlayerCount = gofakeit.Number(2, 8)
	configuration.ZoneConfiguration.PlayerZoneCastles = 2
	configuration.MatchPlayerCastleFactions = true
	templateGenerator := template_generator.NewTemplateGenerator(configuration)

	// Act
	actual := templateGenerator.Generate()

	// Assert
	spawnZones := linq.FromSlice(actual.Variants[0].Zones).
		Where(func(x entities.Zone) bool { return strings.HasPrefix(x.Name, "Spawn-") }).
		ToSlice()
	for _, zone := range spawnZones {
		castle := zone.MainObjects[1]
		if assert.NotNil(t, castle.Faction) {
			assert.Equal(t, "Match", castle.Faction.Type)
		}
	}
}

func TestWhenMatchPlayerCastleFactionsDisabled_SetsRandomFactionOnPlayerCastles(t *testing.T) {
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyDefault
	configuration.PlayerCount = gofakeit.Number(2, 8)
	configuration.ZoneConfiguration.PlayerZoneCastles = 2
	configuration.MatchPlayerCastleFactions = false
	templateGenerator := template_generator.NewTemplateGenerator(configuration)

	// Act
	actual := templateGenerator.Generate()

	// Assert
	spawnZones := linq.FromSlice(actual.Variants[0].Zones).
		Where(func(x entities.Zone) bool { return strings.HasPrefix(x.Name, "Spawn-") }).
		ToSlice()
	for _, zone := range spawnZones {
		castle := zone.MainObjects[1]
		if assert.NotNil(t, castle.Faction) {
			assert.Equal(t, "Random", castle.Faction.Type)
		}
	}
}

// ── Generate: city hold / lost city ──────────────────────────────────

func TestWhenCityHoldEnabled_MarksHoldCityWinConditionObject(t *testing.T) {
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
	templateGenerator := template_generator.NewTemplateGenerator(configuration)

	// Act
	actual := templateGenerator.Generate()

	// Assert
	assert.True(t, actual.GameRules.WinConditions.CityHold)
	holdCityZones := linq.FromSlice(actual.Variants[0].Zones).
		Where(func(zone entities.Zone) bool {
			return linq.FromSlice(zone.MainObjects).
				Where(func(object entities.MainObject) bool { return object.HoldCityWinCon }).
				Any()
		}).
		ToSlice()
	assert.NotEmpty(t, holdCityZones)
}

func TestWhenVictoryConditionFiveProvided_EnablesCityHoldWinCondition(t *testing.T) {
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.GameEndConditions = &config.GameEndConditions{
		VictoryCondition: "win_condition_5",
		LostStartCityDay: 3,
		CityHoldDays:     6,
	}
	templateGenerator := template_generator.NewTemplateGenerator(configuration)

	// Act
	actual := templateGenerator.Generate()

	// Assert
	assert.True(t, actual.GameRules.WinConditions.CityHold)
}

func TestWhenCityHoldEnabledWithHubAndSpokeTopology_MarksHubAsHoldCity(t *testing.T) {
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
	templateGenerator := template_generator.NewTemplateGenerator(configuration)

	// Act
	actual := templateGenerator.Generate()

	// Assert
	hubHoldsCity := linq.FromSlice(actual.Variants[0].Zones).
		Where(func(zone entities.Zone) bool { return zone.Name == "Hub" }).
		Where(func(zone entities.Zone) bool {
			return linq.FromSlice(zone.MainObjects).
				Where(func(object entities.MainObject) bool { return object.HoldCityWinCon }).
				Any()
		}).
		Any()
	assert.True(t, hubHoldsCity)
}

func TestWhenVictoryConditionThreeProvided_EnablesLostStartCityWinCondition(t *testing.T) {
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.GameEndConditions = &config.GameEndConditions{
		VictoryCondition: "win_condition_3",
		LostStartCityDay: 5,
		CityHoldDays:     6,
	}
	templateGenerator := template_generator.NewTemplateGenerator(configuration)

	// Act
	actual := templateGenerator.Generate()

	// Assert
	assert.True(t, actual.GameRules.WinConditions.LostStartCity)
}

// ── Generate: gladiator arena ────────────────────────────────────────

func TestWhenGladiatorArenaEnabled_PropagatesGladiatorRules(t *testing.T) {
	// Arrange
	expectedDaysDelayStart := gofakeit.Number(1, 30)
	expectedCountDay := gofakeit.Number(1, 7)
	configuration := config.NewGeneratorConfig()
	configuration.GladiatorArenaRules = &config.GladiatorArenaRules{
		Enabled:        true,
		DaysDelayStart: expectedDaysDelayStart,
		CountDay:       expectedCountDay,
	}
	templateGenerator := template_generator.NewTemplateGenerator(configuration)

	// Act
	actual := templateGenerator.Generate()

	// Assert
	winConditions := actual.GameRules.WinConditions
	assert.True(t, winConditions.GladiatorArena)
	assert.Equal(t, expectedDaysDelayStart, winConditions.GladiatorArenaDaysDelayStart)
	assert.Equal(t, expectedCountDay, winConditions.GladiatorArenaCountDay)
	assert.Equal(t, "StartHero", winConditions.ChampionSelectRule)
}

func TestWhenVictoryConditionFourProvided_EnablesGladiatorArenaWinCondition(t *testing.T) {
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.GameEndConditions = &config.GameEndConditions{
		VictoryCondition: "win_condition_4",
		LostStartCityDay: 3,
		CityHoldDays:     6,
	}
	templateGenerator := template_generator.NewTemplateGenerator(configuration)

	// Act
	actual := templateGenerator.Generate()

	// Assert
	assert.True(t, actual.GameRules.WinConditions.GladiatorArena)
}

// ── Generate: tournament ─────────────────────────────────────────────

func TestWhenTournamentEnabled_FillsRoundScheduleFromPointsToWin(t *testing.T) {
	// Arrange
	expectedPointsToWin := gofakeit.Number(1, 5)
	configuration := config.NewGeneratorConfig()
	configuration.TournamentRules = &config.TournamentRules{
		Enabled:            true,
		FirstTournamentDay: 10,
		Interval:           5,
		PointsToWin:        expectedPointsToWin,
		SaveArmy:           true,
	}
	templateGenerator := template_generator.NewTemplateGenerator(configuration)

	// Act
	actual := templateGenerator.Generate()

	// Assert
	winConditions := actual.GameRules.WinConditions
	expectedRoundCount := expectedPointsToWin*2 - 1
	assert.True(t, winConditions.Tournament)
	assert.Equal(t, expectedPointsToWin, winConditions.TournamentPointsToWin)
	assert.Len(t, winConditions.TournamentAnnounceDays, expectedRoundCount)
	assert.Len(t, winConditions.TournamentDays, expectedRoundCount)
}

func TestWhenVictoryConditionSixProvided_EnablesTournamentWinCondition(t *testing.T) {
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.GameEndConditions = &config.GameEndConditions{
		VictoryCondition: "win_condition_6",
		LostStartCityDay: 3,
		CityHoldDays:     6,
	}
	templateGenerator := template_generator.NewTemplateGenerator(configuration)

	// Act
	actual := templateGenerator.Generate()

	// Assert
	assert.True(t, actual.GameRules.WinConditions.Tournament)
}

func TestWhenTournamentEnabledWithTwoPlayersAndDefaultTopology_CreatesRingGuardGroups(t *testing.T) {
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyDefault
	configuration.PlayerCount = 2
	configuration.ZoneConfiguration.NeutralZoneCount = gofakeit.Number(2, 6)
	configuration.TournamentRules = &config.TournamentRules{
		Enabled:            true,
		FirstTournamentDay: 14,
		Interval:           7,
		PointsToWin:        2,
	}
	templateGenerator := template_generator.NewTemplateGenerator(configuration)

	// Act
	actual := templateGenerator.Generate()

	// Assert
	hasRingGuardGroup := linq.FromSlice(actual.Variants[0].Connections).
		Where(func(x entities.Connection) bool {
			return strings.HasPrefix(x.GuardMatchGroup, "tourney_ring_guard_")
		}).
		Any()
	assert.True(t, hasRingGuardGroup)
}

func TestWhenTournamentEnabledWithHubAndSpokeTopology_CreatesOneHubPerPlayer(t *testing.T) {
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
	templateGenerator := template_generator.NewTemplateGenerator(configuration)

	// Act
	actual := templateGenerator.Generate()

	// Assert
	perPlayerHubs := linq.FromSlice(actual.Variants[0].Zones).
		Where(func(x entities.Zone) bool { return strings.HasPrefix(x.Name, "Hub-") }).
		ToSlice()
	assert.Equal(t, expectedHubCount, len(perPlayerHubs))
}

func TestWhenTournamentEnabledWithChainTopology_CreatesChainGuardGroups(t *testing.T) {
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
	templateGenerator := template_generator.NewTemplateGenerator(configuration)

	// Act
	actual := templateGenerator.Generate()

	// Assert
	hasChainGuardGroup := linq.FromSlice(actual.Variants[0].Connections).
		Where(func(x entities.Connection) bool {
			return strings.HasPrefix(x.GuardMatchGroup, "tourney_guard_")
		}).
		Any()
	assert.True(t, hasChainGuardGroup)
}

func TestWhenTournamentEnabledWithCirclesTopology_CreatesBalancedGuardGroups(t *testing.T) {
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
	templateGenerator := template_generator.NewTemplateGenerator(configuration)

	// Act
	actual := templateGenerator.Generate()

	// Assert
	hasBalancedGuardGroup := linq.FromSlice(actual.Variants[0].Connections).
		Where(func(x entities.Connection) bool {
			return strings.HasPrefix(x.GuardMatchGroup, "tourney_bal_guard_")
		}).
		Any()
	assert.True(t, hasBalancedGuardGroup)
}

// ── Generate: advanced neutral mix ───────────────────────────────────

func TestWhenAdvancedModeEnabledAndNeutralZonesSelected_SetsExpectedNeutralZoneCount(t *testing.T) {
	// Arrange
	maxCastleCount := 30
	lowNoCastleCount := gofakeit.Number(0, maxCastleCount)
	maxCastleCount -= lowNoCastleCount
	lowCastleCount := gofakeit.Number(0, maxCastleCount)
	maxCastleCount -= lowCastleCount
	mediumNoCastleCount := gofakeit.Number(0, maxCastleCount)
	maxCastleCount -= mediumNoCastleCount
	mediumCastleCount := gofakeit.Number(0, maxCastleCount)
	maxCastleCount -= mediumCastleCount
	highNoCastleCount := gofakeit.Number(0, maxCastleCount)
	maxCastleCount -= highNoCastleCount
	highCastleCount := gofakeit.Number(0, maxCastleCount)
	expectedNeutralCount := lowNoCastleCount + lowCastleCount + mediumNoCastleCount + mediumCastleCount + highNoCastleCount + highCastleCount

	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyDefault
	configuration.ZoneConfiguration.Advanced.Enabled = true
	configuration.ZoneConfiguration.Advanced.NeutralLowNoCastleCount = lowNoCastleCount
	configuration.ZoneConfiguration.Advanced.NeutralLowCastleCount = lowCastleCount
	configuration.ZoneConfiguration.Advanced.NeutralMediumNoCastleCount = mediumNoCastleCount
	configuration.ZoneConfiguration.Advanced.NeutralMediumCastleCount = mediumCastleCount
	configuration.ZoneConfiguration.Advanced.NeutralHighNoCastleCount = highNoCastleCount
	configuration.ZoneConfiguration.Advanced.NeutralHighCastleCount = highCastleCount
	templateGenerator := template_generator.NewTemplateGenerator(configuration)

	// Act
	actual := templateGenerator.Generate()

	// Assert
	actualNeutralCount := len(linq.FromSlice(actual.Variants[0].Zones).
		Where(func(x entities.Zone) bool { return strings.HasPrefix(x.Name, "Neutral-") }).
		ToSlice())
	assert.Equal(t, expectedNeutralCount, actualNeutralCount)
}

func TestWhenAdvancedGuardRandomizationExceedsMaximum_ClampsGuardRandomization(t *testing.T) {
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.ZoneConfiguration.Advanced.Enabled = true
	configuration.ZoneConfiguration.Advanced.GuardRandomization = gofakeit.Float64Range(0.6, 10.0)
	templateGenerator := template_generator.NewTemplateGenerator(configuration)

	// Act
	actual := templateGenerator.Generate()

	// Assert
	for _, zone := range actual.Variants[0].Zones {
		assert.LessOrEqual(t, zone.GuardRandomization, 0.5, "zone %q guard randomization should be clamped", zone.Name)
	}
}

// ── Generate: structural template fields ─────────────────────────────

func TestWhenGenerating_AlwaysProducesExactlyOneVariant(t *testing.T) {
	// Arrange
	configuration := config.NewGeneratorConfig()
	templateGenerator := template_generator.NewTemplateGenerator(configuration)

	// Act
	actual := templateGenerator.Generate()

	// Assert
	assert.Len(t, actual.Variants, 1)
}

func TestWhenGenerating_ProducesFourZoneLayouts(t *testing.T) {
	// Arrange
	configuration := config.NewGeneratorConfig()
	templateGenerator := template_generator.NewTemplateGenerator(configuration)

	// Act
	actual := templateGenerator.Generate()

	// Assert
	assert.Len(t, actual.ZoneLayouts, 4)
}

func TestWhenGenerating_ProducesContentCountLimits(t *testing.T) {
	// Arrange
	configuration := config.NewGeneratorConfig()
	templateGenerator := template_generator.NewTemplateGenerator(configuration)

	// Act
	actual := templateGenerator.Generate()

	// Assert
	assert.NotEmpty(t, actual.ContentCountLimits)
}

func TestWhenChainTopologySelected_IncludesTopologyNameInDescription(t *testing.T) {
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyChain
	configuration.PlayerCount = gofakeit.Number(2, 8)
	configuration.ZoneConfiguration.NeutralZoneCount = gofakeit.Number(0, 5)
	templateGenerator := template_generator.NewTemplateGenerator(configuration)

	// Act
	actual := templateGenerator.Generate()

	// Assert
	assert.Contains(t, actual.Description, "Chain")
}

func TestWhenDescriptionOptionsEnabled_AppendsOptionPhrasesToDescription(t *testing.T) {
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyDefault
	configuration.PlayerCount = gofakeit.Number(2, 8)
	configuration.ZoneConfiguration.NeutralZoneCount = gofakeit.Number(2, 6)
	configuration.NoDirectPlayerConnections = true
	configuration.RandomPortals = true
	configuration.SpawnRemoteFootholds = false
	configuration.GenerateRoads = false
	templateGenerator := template_generator.NewTemplateGenerator(configuration)

	// Act
	actual := templateGenerator.Generate()

	// Assert
	assert.Contains(t, actual.Description, "isolated player starts")
	assert.Contains(t, actual.Description, "random portals")
	assert.Contains(t, actual.Description, "no remote footholds")
	assert.Contains(t, actual.Description, "roads disabled")
}

func TestWhenVictoryConditionProvided_PropagatesToDisplayWinCondition(t *testing.T) {
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.GameEndConditions = &config.GameEndConditions{
		VictoryCondition: "win_condition_2",
		LostStartCityDay: 3,
		CityHoldDays:     6,
	}
	templateGenerator := template_generator.NewTemplateGenerator(configuration)

	// Act
	actual := templateGenerator.Generate()

	// Assert
	assert.Equal(t, "win_condition_2", actual.DisplayWinCondition)
}

func TestWhenGameEndConditionsAreNil_UsesDefaultWinCondition(t *testing.T) {
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.GameEndConditions = nil
	templateGenerator := template_generator.NewTemplateGenerator(configuration)

	// Act
	actual := templateGenerator.Generate()

	// Assert
	assert.Equal(t, "win_condition_1", actual.DisplayWinCondition)
}

// ── Generate: mandatory content groups ───────────────────────────────

func TestWhenGenerating_CreatesMandatoryContentGroupPerPlayerAndNeutralZone(t *testing.T) {
	// Arrange
	playerCount := gofakeit.Number(2, 8)
	neutralZoneCount := gofakeit.Number(1, 6)
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyDefault
	configuration.PlayerCount = playerCount
	configuration.ZoneConfiguration.NeutralZoneCount = neutralZoneCount
	templateGenerator := template_generator.NewTemplateGenerator(configuration)

	// Act
	actual := templateGenerator.Generate()

	// Assert
	playerGroups := linq.FromSlice(actual.MandatoryContent).
		Where(func(x entities.MandatoryContent) bool { return strings.HasPrefix(x.Name, "mandatory_content_side_") }).
		ToSlice()
	neutralGroups := linq.FromSlice(actual.MandatoryContent).
		Where(func(x entities.MandatoryContent) bool { return strings.HasPrefix(x.Name, "mandatory_content_neutral_") }).
		ToSlice()
	assert.Len(t, playerGroups, playerCount)
	assert.Len(t, neutralGroups, neutralZoneCount)
}

// ── Generate: spawn zone main objects ────────────────────────────────

func TestWhenGenerating_PlacesSpawnMainObjectInEachSpawnZone(t *testing.T) {
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyDefault
	configuration.PlayerCount = gofakeit.Number(2, 8)
	templateGenerator := template_generator.NewTemplateGenerator(configuration)

	// Act
	actual := templateGenerator.Generate()

	// Assert
	spawnZones := linq.FromSlice(actual.Variants[0].Zones).
		Where(func(x entities.Zone) bool { return strings.HasPrefix(x.Name, "Spawn-") }).
		ToSlice()
	for _, zone := range spawnZones {
		if assert.NotEmpty(t, zone.MainObjects) {
			assert.Equal(t, "Spawn", zone.MainObjects[0].Type)
		}
	}
}

func TestWhenMultiplePlayerCastlesConfigured_AddsConfiguredCastleCountToEachSpawnZone(t *testing.T) {
	// Arrange
	expectedCastleCount := gofakeit.Number(2, 5)
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyDefault
	configuration.PlayerCount = gofakeit.Number(2, 8)
	configuration.ZoneConfiguration.PlayerZoneCastles = expectedCastleCount
	templateGenerator := template_generator.NewTemplateGenerator(configuration)

	// Act
	actual := templateGenerator.Generate()

	// Assert
	spawnZones := linq.FromSlice(actual.Variants[0].Zones).
		Where(func(x entities.Zone) bool { return strings.HasPrefix(x.Name, "Spawn-") }).
		ToSlice()
	for _, zone := range spawnZones {
		assert.Len(t, zone.MainObjects, expectedCastleCount)
	}
}

// ── Generate: border guards scale with quality ───────────────────────

func TestWhenNeutralZonesAreHighQuality_ProducesHigherBorderGuardValues(t *testing.T) {
	// Arrange
	totalGuardValueFor := func(lowQualityCount, highQualityCount int) int {
		configuration := config.NewGeneratorConfig()
		configuration.Topology = config.TopologyDefault
		configuration.PlayerCount = 2
		configuration.ZoneConfiguration.Advanced.Enabled = true
		configuration.ZoneConfiguration.Advanced.NeutralLowNoCastleCount = lowQualityCount
		configuration.ZoneConfiguration.Advanced.NeutralHighNoCastleCount = highQualityCount
		actual := template_generator.NewTemplateGenerator(configuration).Generate()
		total := 0
		for _, connection := range actual.Variants[0].Connections {
			total += connection.GuardValue
		}
		return total
	}

	// Act
	highQualityTotal := totalGuardValueFor(0, 4)
	lowQualityTotal := totalGuardValueFor(4, 0)

	// Assert
	assert.Greater(t, highQualityTotal, lowQualityTotal)
}

// ── Generate: comprehensive topology smoke ───────────────────────────

func TestWhenGeneratingForEachTopology_ProducesZones(t *testing.T) {
	topologies := []config.MapTopology{
		config.TopologyDefault, config.TopologyChain, config.TopologyHubAndSpoke,
		config.TopologySharedWeb, config.TopologyRandom, config.TopologyCircles,
		config.TopologySquare, config.TopologyGeometric, config.TopologyCross,
		config.TopologyFractal,
	}
	for _, topology := range topologies {
		t.Run(string(topology), func(t *testing.T) {
			// Arrange
			configuration := config.NewGeneratorConfig()
			configuration.Topology = topology
			configuration.PlayerCount = 4
			configuration.ZoneConfiguration.NeutralZoneCount = 4
			templateGenerator := template_generator.NewTemplateGenerator(configuration)

			// Act
			actual := templateGenerator.Generate()

			// Assert
			assert.NotEmpty(t, actual.Variants[0].Zones)
		})
	}
}

// ── Generate: per-topology snapshot assertion ────────────────────────

type expectedConnection struct {
	From, To       string
	ConnectionType string
}

type expectedTemplate struct {
	Name                string
	GameMode            string
	Description         string
	DisplayWinCondition string
	SizeX               int
	SizeZ               int
	VariantCount        int
	ZoneNames           []string // order-independent
	Connections         []expectedConnection
	ZoneLayoutCount     int
	MandatoryCount      int
	ContentCountLimits  int
	ContentPoolCount    int
	ContentListCount    int

	HeroCountMin           int
	HeroCountMax           int
	HeroCountIncrement     int
	HeroHireBan            bool
	FactionLawsExpModifier float64
	AstrologyExpModifier   float64

	WinClassic       bool
	WinDesertion     bool
	WinDesertionDay  int
	WinHeroLighting  bool
	WinLostStartCity int
	WinCityHoldDays  int
}

func TestWhenGeneratingForEachTopology_ReturnsExpectedTemplate(t *testing.T) {
	commonGameRules := struct {
		min, max, inc                            int
		hireBan                                  bool
		factionLawsMod, astrologyMod             float64
		winClassic, winDes, winHeroLight         bool
		winDesDay, winLostStartCity, winCityHold int
	}{4, 8, 1, false, 1.0, 1.0, true, true, true, 3, 3, 6}

	mk := func(name, gameMode, description string, zones []string, conns []expectedConnection, mandatory int) expectedTemplate {
		return expectedTemplate{
			Name:                   name,
			GameMode:               gameMode,
			Description:            description,
			DisplayWinCondition:    "win_condition_1",
			SizeX:                  160,
			SizeZ:                  160,
			VariantCount:           1,
			ZoneNames:              zones,
			Connections:            conns,
			ZoneLayoutCount:        4,
			MandatoryCount:         mandatory,
			ContentCountLimits:     17,
			ContentPoolCount:       0,
			ContentListCount:       0,
			HeroCountMin:           commonGameRules.min,
			HeroCountMax:           commonGameRules.max,
			HeroCountIncrement:     commonGameRules.inc,
			HeroHireBan:            commonGameRules.hireBan,
			FactionLawsExpModifier: commonGameRules.factionLawsMod,
			AstrologyExpModifier:   commonGameRules.astrologyMod,
			WinClassic:             commonGameRules.winClassic,
			WinDesertion:           commonGameRules.winDes,
			WinDesertionDay:        commonGameRules.winDesDay,
			WinHeroLighting:        commonGameRules.winHeroLight,
			WinLostStartCity:       commonGameRules.winLostStartCity,
			WinCityHoldDays:        commonGameRules.winCityHold,
		}
	}

	cases := []struct {
		topology config.MapTopology
		want     expectedTemplate
	}{
		{
			topology: config.TopologyDefault,
			want: mk(
				"Custom Template", "Classic",
				"Generated with Custom Template Editor: Ring layout, no neutral zones, 1 castle per player zone.",
				[]string{"Spawn-A", "Spawn-B"},
				[]expectedConnection{
					{"Spawn-A", "Spawn-B", "Direct"},
					{"Spawn-B", "Spawn-A", "Direct"},
				},
				2,
			),
		},
		{
			topology: config.TopologyChain,
			want: mk(
				"Custom Template", "Classic",
				"Generated with Custom Template Editor: Chain layout, no neutral zones, 1 castle per player zone.",
				[]string{"Spawn-A", "Spawn-B"},
				[]expectedConnection{{"Spawn-A", "Spawn-B", "Direct"}},
				2,
			),
		},
		{
			topology: config.TopologyHubAndSpoke,
			want: mk(
				"Custom Template", "Classic",
				"Generated with Custom Template Editor: Hub layout, no neutral zones, 1 castle per player zone.",
				[]string{"Hub", "Spawn-A", "Spawn-B"},
				[]expectedConnection{
					{"Hub", "Spawn-B", "Direct"},
					{"Hub", "Spawn-B", "Direct"},
					{"Hub", "Spawn-A", "Direct"},
					{"Hub", "Spawn-A", "Direct"},
					{"Spawn-B", "Spawn-A", "Proximity"},
					{"Spawn-A", "Spawn-B", "Proximity"},
				},
				2,
			),
		},
		{
			topology: config.TopologySharedWeb,
			want: mk(
				"Custom Template", "Classic",
				"Generated with Custom Template Editor: Shared Web layout, 1 neutral zone, 1 castle per player zone, 1 castle per neutral zone.",
				[]string{"Neutral-C", "Spawn-A", "Spawn-B"},
				[]expectedConnection{
					{"Spawn-B", "Neutral-C", "Direct"},
					{"Spawn-A", "Neutral-C", "Direct"},
				},
				3,
			),
		},
		{
			topology: config.TopologyRandom,
			want: mk(
				"Custom Template", "Classic",
				"Generated with Custom Template Editor: Random layout, no neutral zones, 1 castle per player zone.",
				[]string{"Spawn-A", "Spawn-B"},
				[]expectedConnection{{"Spawn-B", "Spawn-A", "Direct"}},
				2,
			),
		},
		{
			topology: config.TopologyCircles,
			want: mk(
				"Custom Template", "Classic",
				"Generated with Custom Template Editor: Circles layout, no neutral zones, 1 castle per player zone.",
				[]string{"Spawn-A", "Spawn-B"},
				[]expectedConnection{{"Spawn-A", "Spawn-B", "Direct"}},
				2,
			),
		},
	}

	for _, testCase := range cases {
		t.Run(string(testCase.topology), func(t *testing.T) {
			// Arrange
			configuration := config.NewGeneratorConfig()
			configuration.Topology = testCase.topology
			templateGenerator := template_generator.NewTemplateGenerator(configuration)

			// Act
			actual := templateGenerator.Generate()

			// Assert
			assertTemplateMatches(t, testCase.want, actual)
		})
	}
}

func assertTemplateMatches(t *testing.T, expected expectedTemplate, actual *entities.RmgTemplate) {
	t.Helper()

	// Top-level scalar fields.
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, expected.GameMode, actual.GameMode)
	assert.Equal(t, expected.Description, actual.Description)
	assert.Equal(t, expected.DisplayWinCondition, actual.DisplayWinCondition)
	assert.Equal(t, expected.SizeX, actual.SizeX)
	assert.Equal(t, expected.SizeZ, actual.SizeZ)

	// Variants / zones / connections.
	assert.Len(t, actual.Variants, expected.VariantCount)
	generatedVariant := actual.Variants[0]

	actualZoneNames := make([]string, 0, len(generatedVariant.Zones))
	for _, zone := range generatedVariant.Zones {
		actualZoneNames = append(actualZoneNames, zone.Name)
	}
	assert.ElementsMatch(t, expected.ZoneNames, actualZoneNames)

	actualConnections := make([]expectedConnection, 0, len(generatedVariant.Connections))
	for _, connection := range generatedVariant.Connections {
		actualConnections = append(actualConnections, normalizeConnection(connection.From, connection.To, connection.ConnectionType))
	}
	expectedConnections := make([]expectedConnection, 0, len(expected.Connections))
	for _, connection := range expected.Connections {
		expectedConnections = append(expectedConnections, normalizeConnection(connection.From, connection.To, connection.ConnectionType))
	}
	assert.ElementsMatch(t, expectedConnections, actualConnections)

	// Collection counts.
	assert.Len(t, actual.ZoneLayouts, expected.ZoneLayoutCount)
	assert.Len(t, actual.MandatoryContent, expected.MandatoryCount)
	assert.Len(t, actual.ContentCountLimits, expected.ContentCountLimits)
	assert.Len(t, actual.ContentPools, expected.ContentPoolCount)
	assert.Len(t, actual.ContentLists, expected.ContentListCount)

	// Game rules.
	gameRules := actual.GameRules
	assert.Equal(t, expected.HeroCountMin, gameRules.HeroCountMin)
	assert.Equal(t, expected.HeroCountMax, gameRules.HeroCountMax)
	assert.Equal(t, expected.HeroCountIncrement, gameRules.HeroCountIncrement)
	assert.Equal(t, expected.HeroHireBan, gameRules.HeroHireBan)
	assert.Equal(t, expected.FactionLawsExpModifier, gameRules.FactionLawsExpModifier)
	assert.Equal(t, expected.AstrologyExpModifier, gameRules.AstrologyExpModifier)

	// Win conditions.
	winConditions := gameRules.WinConditions
	assert.Equal(t, expected.WinClassic, winConditions.Classic)
	assert.Equal(t, expected.WinDesertion, winConditions.Desertion)
	assert.Equal(t, expected.WinDesertionDay, winConditions.DesertionDay)
	assert.Equal(t, expected.WinHeroLighting, winConditions.HeroLighting)
	assert.Equal(t, expected.WinLostStartCity, winConditions.LostStartCityDay)
	assert.Equal(t, expected.WinCityHoldDays, winConditions.CityHoldDays)
}

func normalizeConnection(from, to, connectionType string) expectedConnection {
	// Connections are conceptually undirected — normalize endpoints so callers
	// don't have to anticipate the (non-deterministic) emission order.
	if from > to {
		from, to = to, from
	}
	return expectedConnection{From: from, To: to, ConnectionType: connectionType}
}

// ── Generate: connection endpoints resolve ───────────────────────────

func TestWhenGenerating_AllConnectionEndpointsReferenceKnownZones(t *testing.T) {
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyDefault
	configuration.PlayerCount = gofakeit.Number(2, 8)
	configuration.ZoneConfiguration.NeutralZoneCount = gofakeit.Number(1, 6)
	templateGenerator := template_generator.NewTemplateGenerator(configuration)

	// Act
	actual := templateGenerator.Generate()

	// Assert
	knownZoneNames := make(map[string]bool)
	for _, zone := range actual.Variants[0].Zones {
		knownZoneNames[zone.Name] = true
	}
	for _, connection := range actual.Variants[0].Connections {
		assert.True(t, knownZoneNames[connection.From], "connection %q references unknown zone %q", connection.Name, connection.From)
		assert.True(t, knownZoneNames[connection.To], "connection %q references unknown zone %q", connection.Name, connection.To)
	}
}
