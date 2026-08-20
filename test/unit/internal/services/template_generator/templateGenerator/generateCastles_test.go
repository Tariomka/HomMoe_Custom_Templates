package templateGenerator_test

import (
	"slices"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newAbandonedOutpostConfiguration builds a deterministic two-player ring
// configuration with one low- and one medium-tier neutral castle zone, so
// abandoned-outpost behaviour can be compared with and without the option.
func newAbandonedOutpostConfiguration(spawnOutposts bool) *config.GeneratorConfig {
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyRing
	configuration.PlayerCount = 2
	configuration.ZoneConfiguration.Advanced.Enabled = true
	configuration.ZoneConfiguration.Advanced.NeutralLowCastleCount = 1
	configuration.ZoneConfiguration.Advanced.NeutralMediumCastleCount = 1
	configuration.ZoneConfiguration.Advanced.NeutralLowCastlesPerZone = 1
	configuration.ZoneConfiguration.Advanced.NeutralMediumCastlesPerZone = 1
	configuration.ZoneConfiguration.SpawnAbandonedOutposts = spawnOutposts
	return configuration
}

// countNeutralMainObjectsOfType counts main objects of the given type across
// all neutral zones of the template's first variant.
func countNeutralMainObjectsOfType(generated *entities.RmgTemplate, objectType string) int {
	count := 0
	for _, zone := range zonesWithPrefix(generated, "Neutral-") {
		for _, mainObject := range zone.MainObjects {
			if mainObject.Type == objectType {
				count++
			}
		}
	}
	return count
}

func TestWhenAbandonedOutpostsDisabled_AddsNoAbandonedOutpostMainObjects(t *testing.T) {
	t.Parallel()
	// Arrange
	generator := test_helpers.NewTemplateGenerator(newAbandonedOutpostConfiguration(false))

	// Act
	actual, _ := generator.Generate()

	// Assert
	assert.Zero(t, countNeutralMainObjectsOfType(actual, "AbandonedOutpost"))
}

func TestWhenAbandonedOutpostsEnabled_AddsAbandonedOutpostMainObjects(t *testing.T) {
	t.Parallel()
	// Arrange
	generator := test_helpers.NewTemplateGenerator(newAbandonedOutpostConfiguration(true))

	// Act
	actual, _ := generator.Generate()

	// Assert
	assert.Positive(t, countNeutralMainObjectsOfType(actual, "AbandonedOutpost"))
}

func TestWhenAbandonedOutpostsEnabled_KeepsNeutralCityCount(t *testing.T) {
	t.Parallel()
	// Arrange
	baseline, _ := test_helpers.NewTemplateGenerator(newAbandonedOutpostConfiguration(false)).Generate()
	baselineCityCount := countNeutralMainObjectsOfType(baseline, "City")
	require.Positive(t, baselineCityCount, "baseline must produce neutral cities to compare against")
	generator := test_helpers.NewTemplateGenerator(newAbandonedOutpostConfiguration(true))

	// Act
	actual, _ := generator.Generate()

	// Assert
	assert.Equal(t, baselineCityCount, countNeutralMainObjectsOfType(actual, "City"))
}

// newPlayerOwnedCastlesConfiguration builds a deterministic two-player ring
// configuration with one unclaimed extra castle and the given number of
// pre-owned castles per spawn zone.
func newPlayerOwnedCastlesConfiguration(ownedPerZone int) *config.GeneratorConfig {
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyRing
	configuration.PlayerCount = 2
	configuration.ZoneConfiguration.PlayerZoneCastles = 1
	configuration.ZoneConfiguration.PlayerOwnedCastles = ownedPerZone
	return configuration
}

// countZoneCitiesWhere counts City main objects of the zone matching the
// given predicate.
func countZoneCitiesWhere(zone entities.Zone, predicate func(entities.MainObject) bool) int {
	count := 0
	for _, mainObject := range zone.MainObjects {
		if mainObject.Type == "City" && predicate(mainObject) {
			count++
		}
	}
	return count
}

func TestWhenPlayerOwnedCastlesConfigured_AddsOwnedCityPerCountInEachSpawnZone(t *testing.T) {
	t.Parallel()
	// Arrange
	const ownedPerZone = 2
	generator := test_helpers.NewTemplateGenerator(newPlayerOwnedCastlesConfiguration(ownedPerZone))

	// Act
	actual, _ := generator.Generate()

	// Assert
	var ownedCounts []int
	for _, zone := range zonesWithPrefix(actual, "Spawn-") {
		ownedCounts = append(ownedCounts, countZoneCitiesWhere(zone,
			func(mainObject entities.MainObject) bool { return mainObject.Owner != "" }))
	}
	assert.Equal(t, []int{ownedPerZone, ownedPerZone}, ownedCounts)
}

func TestWhenPlayerOwnedCastlesConfigured_AssignsSpawnPlayerAsOwner(t *testing.T) {
	t.Parallel()
	// Arrange
	generator := test_helpers.NewTemplateGenerator(newPlayerOwnedCastlesConfiguration(gofakeit.Number(1, 5)))

	// Act
	actual, _ := generator.Generate()

	// Assert
	var ownerMismatches []string
	for _, zone := range zonesWithPrefix(actual, "Spawn-") {
		spawnPlayer := zone.MainObjects[0].Spawn
		for _, mainObject := range zone.MainObjects {
			if mainObject.Type == "City" && mainObject.Owner != "" && mainObject.Owner != spawnPlayer {
				ownerMismatches = append(ownerMismatches, zone.Name+": owner "+mainObject.Owner+", want "+spawnPlayer)
			}
		}
	}
	assert.Empty(t, ownerMismatches)
}

func TestWhenPlayerOwnedCastlesConfigured_KeepsConfiguredUnclaimedCastleCount(t *testing.T) {
	t.Parallel()
	// Arrange
	generator := test_helpers.NewTemplateGenerator(newPlayerOwnedCastlesConfiguration(gofakeit.Number(1, 5)))

	// Act
	actual, _ := generator.Generate()

	// Assert
	var unclaimedCounts []int
	for _, zone := range zonesWithPrefix(actual, "Spawn-") {
		unclaimedCounts = append(unclaimedCounts, countZoneCitiesWhere(zone,
			func(mainObject entities.MainObject) bool { return mainObject.Owner == "" }))
	}
	assert.Equal(t, []int{1, 1}, unclaimedCounts)
}

func TestWhenPlayerOwnedCastlesConfigured_UnclaimedCastlesKeepGuards(t *testing.T) {
	t.Parallel()
	// Arrange
	generator := test_helpers.NewTemplateGenerator(newPlayerOwnedCastlesConfiguration(gofakeit.Number(1, 5)))

	// Act
	actual, _ := generator.Generate()

	// Assert
	var guardViolations []string
	for _, zone := range zonesWithPrefix(actual, "Spawn-") {
		for _, mainObject := range zone.MainObjects {
			if mainObject.Type == "City" && mainObject.Owner == "" && mainObject.RemoveGuardIfHasOwner {
				guardViolations = append(guardViolations, zone.Name+": unclaimed castle should keep its guard")
			}
		}
	}
	assert.Empty(t, guardViolations)
}

func TestWhenPlayerZoneCastlesConfigured_CreatesSpawnPlusConfiguredCastleMainObjects(t *testing.T) {
	t.Parallel()
	// Arrange
	extraCastles := gofakeit.Number(1, 5)
	playerCount := gofakeit.Number(2, 8)
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyRing
	configuration.PlayerCount = playerCount
	configuration.ZoneConfiguration.PlayerZoneCastles = extraCastles
	generator := test_helpers.NewTemplateGenerator(configuration)

	// Act
	actual, _ := generator.Generate()

	// Assert
	var mainObjectCounts []int
	for _, zone := range zonesWithPrefix(actual, "Spawn-") {
		mainObjectCounts = append(mainObjectCounts, len(zone.MainObjects))
	}
	assert.Equal(t, slices.Repeat([]int{1 + extraCastles}, playerCount), mainObjectCounts)
}

func TestWhenGenerating_AssignsPlayerToEachSpawnMainObject(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyRing
	configuration.PlayerCount = gofakeit.Number(2, 8)
	generator := test_helpers.NewTemplateGenerator(configuration)

	// Act
	actual, _ := generator.Generate()

	// Assert
	var spawnAssignments []string
	for _, zone := range zonesWithPrefix(actual, "Spawn-") {
		spawnAssignments = append(spawnAssignments, zone.MainObjects[0].Spawn)
	}
	assert.NotContains(t, spawnAssignments, "")
}
