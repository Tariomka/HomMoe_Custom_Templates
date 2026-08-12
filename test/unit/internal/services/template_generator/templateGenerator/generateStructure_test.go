package templateGenerator_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenRingTopologySelected_CreatesConnectionPerZone(t *testing.T) {
	t.Parallel()
	// Arrange
	playerCount := gofakeit.Number(3, 8)
	neutralZoneCount := gofakeit.Number(1, 6)
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyRing
	configuration.PlayerCount = playerCount
	configuration.ZoneConfiguration.NeutralZoneCount = neutralZoneCount
	generator := test_helpers.NewTemplateGenerator(configuration)

	// Act
	actual, _ := generator.Generate()

	// Assert
	assert.Len(t, actual.Variants[0].Connections, playerCount+neutralZoneCount)
}

func TestWhenRingTopologyWithEightZones_SetsOrientationAngleStepToFortyFiveDegrees(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyRing
	configuration.PlayerCount = 4
	configuration.ZoneConfiguration.NeutralZoneCount = 4
	configuration.ShufflePlayerZones = false // Deterministic zero-angle zone.
	generator := test_helpers.NewTemplateGenerator(configuration)

	// Act
	actual, _ := generator.Generate()

	// Assert
	expectedOrientation := entities.Orientation{
		ZeroAngleZone:        "Spawn-A",
		BaseAngleMin:         45,
		BaseAngleMax:         45,
		RandomAngleAmplitude: 360,
		RandomAngleStep:      45, // 8 zones -> 360 / 8.
	}
	assert.Equal(t, expectedOrientation, actual.Variants[0].Orientation)
}

func TestWhenTopologySelected_IncludesTopologyNameInDescription(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name           string
		topology       config.MapTopology
		expectedPhrase string
	}{
		{"WhenRingTopologySelected_IncludesRingInDescription", config.TopologyRing, "Ring"},
		{"WhenHubAndSpokeTopologySelected_IncludesHubInDescription", config.TopologyHubAndSpoke, "Hub"},
		{"WhenSharedWebTopologySelected_IncludesSharedWebInDescription", config.TopologySharedWeb, "Shared Web"},
		{"WhenRandomTopologySelected_IncludesRandomInDescription", config.TopologyRandom, "Random"},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			// Arrange
			configuration := config.NewGeneratorConfig()
			configuration.Topology = testCase.topology
			configuration.PlayerCount = 2
			configuration.ZoneConfiguration.NeutralZoneCount = 1
			generator := test_helpers.NewTemplateGenerator(configuration)

			// Act
			actual, _ := generator.Generate()

			// Assert
			assert.Contains(t, actual.Description, testCase.expectedPhrase)
		})
	}
}

// newCirclesMixedNeutralConfiguration builds a four-player circles
// configuration with two castle-free neutral zones of every quality tier.
func newCirclesMixedNeutralConfiguration() *config.GeneratorConfig {
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyCircles
	configuration.PlayerCount = 4
	configuration.ZoneConfiguration.NeutralZoneCount = 0
	configuration.ZoneConfiguration.Advanced.Enabled = true
	configuration.ZoneConfiguration.Advanced.NeutralLowNoCastleCount = 2
	configuration.ZoneConfiguration.Advanced.NeutralMediumNoCastleCount = 2
	configuration.ZoneConfiguration.Advanced.NeutralHighNoCastleCount = 2
	return configuration
}

func TestWhenCirclesTopologyWithMixedNeutralTiers_CreatesZoneForEveryPlannedZone(t *testing.T) {
	t.Parallel()
	// Arrange
	generator := test_helpers.NewTemplateGenerator(newCirclesMixedNeutralConfiguration())

	// Act
	actual, _ := generator.Generate()

	// Assert
	assert.GreaterOrEqual(t, len(actual.Variants[0].Zones), 10, "4 player + 6 neutral zones expected")
}

func TestWhenCirclesTopologyWithMixedNeutralTiers_CreatesConnections(t *testing.T) {
	t.Parallel()
	// Arrange
	generator := test_helpers.NewTemplateGenerator(newCirclesMixedNeutralConfiguration())

	// Act
	actual, _ := generator.Generate()

	// Assert
	assert.NotEmpty(t, actual.Variants[0].Connections)
}

// sumConnectionGuardValues totals the guard values of every connection in the
// template's first variant.
func sumConnectionGuardValues(generated *entities.RmgTemplate) int {
	total := 0
	for _, connection := range generated.Variants[0].Connections {
		total += connection.GuardValue
	}
	return total
}

func TestWhenNeutralZonesAreHighQuality_ProducesStrongerBorderGuardsThanLowQuality(t *testing.T) {
	t.Parallel()
	// Arrange
	newQualityConfiguration := func(highCount, lowCount int) *config.GeneratorConfig {
		configuration := config.NewGeneratorConfig()
		configuration.Topology = config.TopologyRing
		configuration.PlayerCount = 2
		configuration.ZoneConfiguration.NeutralZoneCount = 0
		configuration.ZoneConfiguration.Advanced.Enabled = true
		configuration.ZoneConfiguration.Advanced.NeutralHighNoCastleCount = highCount
		configuration.ZoneConfiguration.Advanced.NeutralLowNoCastleCount = lowCount
		return configuration
	}
	highQualityGenerator := test_helpers.NewTemplateGenerator(newQualityConfiguration(4, 0))
	lowQualityGenerator := test_helpers.NewTemplateGenerator(newQualityConfiguration(0, 4))

	// Act
	highQualityTemplate, _ := highQualityGenerator.Generate()
	lowQualityTemplate, _ := lowQualityGenerator.Generate()
	highQualityTotal := sumConnectionGuardValues(highQualityTemplate)
	lowQualityTotal := sumConnectionGuardValues(lowQualityTemplate)

	// Assert
	assert.Greater(t, highQualityTotal, lowQualityTotal)
}

func TestWhenCityHoldEnabledWithMixedNeutralTiers_MarksExactlyOneHoldCityMainObject(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyHubAndSpoke
	configuration.PlayerCount = 2
	configuration.GameEndConditions = &config.GameEndConditions{
		VictoryCondition: "win_condition_5",
		CityHoldDays:     6,
		LostStartCityDay: 3,
	}
	configuration.ZoneConfiguration.Advanced.Enabled = true
	configuration.ZoneConfiguration.Advanced.NeutralLowNoCastleCount = 2
	configuration.ZoneConfiguration.Advanced.NeutralMediumCastleCount = 2
	configuration.ZoneConfiguration.Advanced.NeutralMediumCastlesPerZone = 1
	configuration.ZoneConfiguration.Advanced.NeutralHighCastleCount = 2
	configuration.ZoneConfiguration.Advanced.NeutralHighCastlesPerZone = 1
	generator := test_helpers.NewTemplateGenerator(configuration)

	// Act
	actual, _ := generator.Generate()

	// Assert
	var holdCityObjects []entities.MainObject
	for _, zone := range actual.Variants[0].Zones {
		for _, mainObject := range zone.MainObjects {
			if mainObject.HoldCityWinCon {
				holdCityObjects = append(holdCityObjects, mainObject)
			}
		}
	}
	assert.Len(t, holdCityObjects, 1)
}
