package connectionEditor_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/connection_editor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeNeutralZone(label string, quality models.NeutralZoneQuality, castles int) entities.Zone {
	return connection_editor.NewDefaultNeutralZone(label, quality, castles, false, defaultTuning())
}

func makeSpawnZone(label, playerName string, extraCastles int) entities.Zone {
	mainObjects := []entities.MainObject{{Type: "Spawn", Spawn: playerName}}
	for range extraCastles {
		mainObjects = append(mainObjects, entities.MainObject{Type: "City"})
	}
	return entities.Zone{Name: "Spawn-" + label, MainObjects: mainObjects}
}

// ════════════════════════════════════════════════════════════════════════
// EditorStateDto.DiffCastleSettings
// ════════════════════════════════════════════════════════════════════════

func TestDiffCastleSettings_NothingChanged_ReportsNothing(t *testing.T) {
	previous := dtos.NewDefaultEditorStateDto()
	current := previous

	assert.False(t, previous.DiffCastleSettings(&current).Any())
}

func TestDiffCastleSettings_SimpleMode_TracksNeutralZoneCastlesOnly(t *testing.T) {
	previous := dtos.NewDefaultEditorStateDto()
	current := previous
	current.NeutralZoneCastles = 3
	current.NeutralHighCastlesPerZone = 4 // advanced-only option, must be ignored

	changes := previous.DiffCastleSettings(&current)

	assert.True(t, changes.NeutralSimple)
	assert.False(t, changes.NeutralLow || changes.NeutralMedium || changes.NeutralHigh)
	assert.False(t, changes.PlayerCastles || changes.Hub)
}

func TestDiffCastleSettings_AdvancedMode_TracksPerTierCounts(t *testing.T) {
	previous := dtos.NewDefaultEditorStateDto()
	previous.AdvancedMode = true
	current := previous
	current.NeutralHighCastlesPerZone = 4
	current.NeutralZoneCastles = 3 // simple-only option, must be ignored

	changes := previous.DiffCastleSettings(&current)

	assert.True(t, changes.NeutralHigh)
	assert.False(t, changes.NeutralSimple)
	assert.False(t, changes.NeutralLow || changes.NeutralMedium)
}

func TestDiffCastleSettings_PlayerAndHubCounts(t *testing.T) {
	previous := dtos.NewDefaultEditorStateDto()
	current := previous
	current.PlayerOwnedCastles = 2
	current.HubZoneCastles = 3

	changes := previous.DiffCastleSettings(&current)

	assert.True(t, changes.PlayerCastles)
	assert.True(t, changes.Hub)
}

// ════════════════════════════════════════════════════════════════════════
// SetNeutralZoneCastleCount
// ════════════════════════════════════════════════════════════════════════

func TestSetNeutralZoneCastleCount_ChangesOnlyCastles(t *testing.T) {
	zone := makeNeutralZone("G", models.QualityHigh, 1)
	originalGuardMultiplier := zone.GuardMultiplier
	originalPool := zone.GuardedContentPool
	originalGuardedValue := zone.GuardedContentValue
	zone.Size = 1.7 // pretend manual size edit

	connection_editor.SetNeutralZoneCastleCount(&zone, 3, defaultTuning())

	assert.Equal(t, 3, connection_editor.CountZoneCastles(zone))
	assert.Equal(t, models.QualityHigh, connection_editor.QualityOfZone(zone))
	assert.Equal(t, originalGuardMultiplier, zone.GuardMultiplier, "guard multiplier must not be re-profiled")
	assert.Equal(t, originalPool, zone.GuardedContentPool, "content pools must not be re-profiled")
	assert.Equal(t, originalGuardedValue, zone.GuardedContentValue, "content values must not be re-profiled")
	assert.Equal(t, 1.7, zone.Size)
}

func TestSetNeutralZoneCastleCount_PreservesNonCastleMainObjects(t *testing.T) {
	zone := makeNeutralZone("G", models.QualityMedium, 1)
	zone.MainObjects = append(zone.MainObjects, entities.MainObject{Type: "AbandonedOutpost"})

	connection_editor.SetNeutralZoneCastleCount(&zone, 2, defaultTuning())

	assert.Equal(t, 2, connection_editor.CountZoneCastles(zone))
	outposts := 0
	for _, mainObject := range zone.MainObjects {
		if mainObject.Type == "AbandonedOutpost" {
			outposts++
		}
	}
	assert.Equal(t, 1, outposts, "abandoned outposts must survive a castle rebuild")
}

func TestSetNeutralZoneCastleCount_PreservesHoldCity(t *testing.T) {
	zone := makeNeutralZone("G", models.QualityHigh, 1)
	zone.MainObjects[0].HoldCityWinCon = true

	connection_editor.SetNeutralZoneCastleCount(&zone, 2, defaultTuning())

	require.Equal(t, 2, connection_editor.CountZoneCastles(zone))
	assert.True(t, zone.MainObjects[0].HoldCityWinCon, "hold-city win condition was lost by the rebuild")
}

// ════════════════════════════════════════════════════════════════════════
// ApplyCastleSettingChanges
// ════════════════════════════════════════════════════════════════════════

func TestApplyCastleSettingChanges_NoChanges_LeavesZonesUntouched(t *testing.T) {
	configuration := config.NewGeneratorConfig()
	zones := []entities.Zone{makeNeutralZone("G", models.QualityMedium, 1)}

	connection_editor.ApplyCastleSettingChanges(zones, dtos.CastleSettingChanges{}, configuration)

	assert.Equal(t, 1, connection_editor.CountZoneCastles(zones[0]))
}

func TestApplyCastleSettingChanges_SimpleMode_UpdatesEveryNeutralZone(t *testing.T) {
	configuration := config.NewGeneratorConfig()
	configuration.ZoneConfiguration.NeutralZoneCastles = 2
	zones := []entities.Zone{
		makeNeutralZone("G", models.QualityMedium, 1),
		makeNeutralZone("H", models.QualityHigh, 0), // castle-less zones update too in simple mode
		makeSpawnZone("A", "Player1", 1),
	}
	spawnCastlesBefore := len(zones[2].MainObjects)

	connection_editor.ApplyCastleSettingChanges(
		zones, dtos.CastleSettingChanges{NeutralSimple: true}, configuration)

	assert.Equal(t, 2, connection_editor.CountZoneCastles(zones[0]))
	assert.Equal(t, 2, connection_editor.CountZoneCastles(zones[1]))
	assert.Len(t, zones[2].MainObjects, spawnCastlesBefore, "spawn zones must not react to a neutral option")
}

func TestApplyCastleSettingChanges_AdvancedMode_MatchesManualQuality(t *testing.T) {
	// The user re-tiered a zone to High in the manual editor; the High
	// castles-per-zone option must drive it by its CURRENT quality, not the
	// quality the generator originally planned.
	configuration := config.NewGeneratorConfig()
	configuration.ZoneConfiguration.Advanced.NeutralHighCastlesPerZone = 3
	zones := []entities.Zone{
		makeNeutralZone("G", models.QualityHigh, 1), // was Low, manually re-tiered
		makeNeutralZone("H", models.QualityLow, 1),  // quality does not match -> untouched
		makeNeutralZone("I", models.QualityHigh, 0), // castle-less -> keeps its no-castle plan
	}

	connection_editor.ApplyCastleSettingChanges(
		zones, dtos.CastleSettingChanges{NeutralHigh: true}, configuration)

	assert.Equal(t, 3, connection_editor.CountZoneCastles(zones[0]))
	assert.Equal(t, 1, connection_editor.CountZoneCastles(zones[1]))
	assert.Equal(t, 0, connection_editor.CountZoneCastles(zones[2]))
}

func TestApplyCastleSettingChanges_PlayerCastles_RebuildsSpawnZone(t *testing.T) {
	configuration := config.NewGeneratorConfig()
	configuration.ZoneConfiguration.PlayerZoneCastles = 2
	configuration.ZoneConfiguration.PlayerOwnedCastles = 1
	zones := []entities.Zone{
		makeSpawnZone("A", "Player1", 0),
		makeNeutralZone("G", models.QualityMedium, 1),
	}

	connection_editor.ApplyCastleSettingChanges(
		zones, dtos.CastleSettingChanges{PlayerCastles: true}, configuration)

	spawn := zones[0]
	// Spawn castle + 1 owned + 2 unclaimed.
	require.Len(t, spawn.MainObjects, 4)
	assert.Equal(t, "Spawn", spawn.MainObjects[0].Type, "the spawn castle must stay the primary main object")
	assert.Equal(t, "Player1", spawn.MainObjects[0].Spawn, "the player assignment was lost")
	owned := 0
	for _, mainObject := range spawn.MainObjects[1:] {
		assert.Equal(t, "City", mainObject.Type)
		if mainObject.Owner == "Player1" {
			owned++
		}
	}
	assert.Equal(t, 1, owned, "expected exactly one player-owned extra castle")
	assert.Equal(t, 1, connection_editor.CountZoneCastles(zones[1]), "neutral zones must not react to a player option")
}

func TestApplyCastleSettingChanges_Hub_RebuildsHubZone(t *testing.T) {
	configuration := config.NewGeneratorConfig()
	configuration.ZoneConfiguration.HubZoneCastles = 3
	zones := []entities.Zone{
		{Name: "Hub", MainObjects: []entities.MainObject{{Type: "City"}}},
	}

	connection_editor.ApplyCastleSettingChanges(
		zones, dtos.CastleSettingChanges{Hub: true}, configuration)

	assert.Equal(t, 3, connection_editor.CountZoneCastles(zones[0]))
}

func TestApplyCastleSettingChanges_RebuildsCastleRoads(t *testing.T) {
	configuration := config.NewGeneratorConfig()
	configuration.ZoneConfiguration.NeutralZoneCastles = 3
	zones := []entities.Zone{makeNeutralZone("G", models.QualityMedium, 1)}

	connection_editor.ApplyCastleSettingChanges(
		zones, dtos.CastleSettingChanges{NeutralSimple: true}, configuration)

	castleRoads := 0
	for _, road := range zones[0].Roads {
		if road.From.Type == "MainObject" && road.To.Type == "MainObject" {
			castleRoads++
		}
	}
	assert.Equal(t, 2, castleRoads, "expected castle roads from the primary castle to both extra castles")
}
