package manualReapplyService_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenNoChangeIsFlagged_LeavesZoneCastlesUntouched(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	zones := []entities.Zone{makeNeutralZone("G", neutral_zone.QualityMedium, 1)}

	// Act
	newManualReapplyService().
		ApplyCastleSettingChanges(zones, editor_state_dto.CastleSettingChanges{}, configuration)

	// Assert
	assert.Equal(t, 1, test_helpers.NewZoneEditorService().CountZoneCastles(zones[0]))
}

// applySimpleModeChange runs the simple-mode neutral castle propagation over a
// castled neutral zone, a castle-less neutral zone and a spawn zone.
func applySimpleModeChange() []entities.Zone {
	configuration := config.NewGeneratorConfig()
	configuration.ZoneConfiguration.NeutralZoneCastles = 2
	zones := []entities.Zone{
		makeNeutralZone("G", neutral_zone.QualityMedium, 1),
		makeNeutralZone("H", neutral_zone.QualityHigh, 0),
		makeSpawnZone("A", "Player1", 1),
	}
	newManualReapplyService().ApplyCastleSettingChanges(
		zones, editor_state_dto.CastleSettingChanges{NeutralSimple: true}, configuration)
	return zones
}

func TestWhenSimpleModeCountChanges_UpdatesCastledNeutralZone(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	zones := applySimpleModeChange()

	// Assert
	assert.Equal(t, 2, test_helpers.NewZoneEditorService().CountZoneCastles(zones[0]))
}

func TestWhenSimpleModeCountChanges_UpdatesCastleLessNeutralZoneToo(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	zones := applySimpleModeChange()

	// Assert
	assert.Equal(t, 2, test_helpers.NewZoneEditorService().CountZoneCastles(zones[1]))
}

func TestWhenSimpleModeCountChanges_LeavesSpawnZoneUntouched(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	zones := applySimpleModeChange()

	// Assert
	assert.Len(t, zones[2].MainObjects, 2, "spawn zones must not react to a neutral option")
}

// applyAdvancedHighChange runs the advanced-mode high-tier castle propagation
// over a manually re-tiered high zone, a low zone and a castle-less high zone.
func applyAdvancedHighChange() []entities.Zone {
	configuration := config.NewGeneratorConfig()
	configuration.ZoneConfiguration.Advanced.NeutralHighCastlesPerZone = 3
	zones := []entities.Zone{
		makeNeutralZone("G", neutral_zone.QualityHigh, 1),
		makeNeutralZone("H", neutral_zone.QualityLow, 1),
		makeNeutralZone("I", neutral_zone.QualityHigh, 0),
	}
	newManualReapplyService().ApplyCastleSettingChanges(
		zones, editor_state_dto.CastleSettingChanges{NeutralHigh: true}, configuration)
	return zones
}

func TestWhenHighTierCountChanges_UpdatesZoneWithMatchingQuality(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	zones := applyAdvancedHighChange()

	// Assert
	assert.Equal(t, 3, test_helpers.NewZoneEditorService().CountZoneCastles(zones[0]))
}

func TestWhenHighTierCountChanges_LeavesOtherQualityZoneUntouched(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	zones := applyAdvancedHighChange()

	// Assert
	assert.Equal(t, 1, test_helpers.NewZoneEditorService().CountZoneCastles(zones[1]))
}

func TestWhenHighTierCountChanges_KeepsCastleLessZoneWithoutCastles(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	zones := applyAdvancedHighChange()

	// Assert
	assert.Equal(t, 0, test_helpers.NewZoneEditorService().CountZoneCastles(zones[2]))
}

// applyPlayerCastleChange runs the player-castle propagation over a spawn zone
// and a neutral zone: one owned extra castle and two unclaimed ones.
func applyPlayerCastleChange() []entities.Zone {
	configuration := config.NewGeneratorConfig()
	configuration.ZoneConfiguration.PlayerZoneCastles = 2
	configuration.ZoneConfiguration.PlayerOwnedCastles = 1
	zones := []entities.Zone{
		makeSpawnZone("A", "Player1", 0),
		makeNeutralZone("G", neutral_zone.QualityMedium, 1),
	}
	newManualReapplyService().ApplyCastleSettingChanges(
		zones, editor_state_dto.CastleSettingChanges{PlayerCastles: true}, configuration)
	return zones
}

func TestWhenPlayerCastlesChange_RebuildsSpawnZoneWithSpawnPlusThreeCastles(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	zones := applyPlayerCastleChange()

	// Assert
	assert.Len(t, zones[0].MainObjects, 4, "expected spawn castle + 1 owned + 2 unclaimed")
}

func TestWhenPlayerCastlesChange_KeepsSpawnCastlePrimary(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	zones := applyPlayerCastleChange()

	// Assert
	require.NotEmpty(t, zones[0].MainObjects)
	assert.Equal(t, "Spawn", zones[0].MainObjects[0].Type,
		"the spawn castle must stay the primary main object")
}

func TestWhenPlayerCastlesChange_KeepsPlayerAssignment(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	zones := applyPlayerCastleChange()

	// Assert
	require.NotEmpty(t, zones[0].MainObjects)
	assert.Equal(t, "Player1", zones[0].MainObjects[0].Spawn, "the player assignment was lost")
}

func TestWhenPlayerCastlesChange_MakesEveryExtraMainObjectACity(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	zones := applyPlayerCastleChange()

	// Assert
	require.Len(t, zones[0].MainObjects, 4)
	extraTypes := []string{
		zones[0].MainObjects[1].Type,
		zones[0].MainObjects[2].Type,
		zones[0].MainObjects[3].Type,
	}
	assert.Equal(t, []string{"City", "City", "City"}, extraTypes)
}

func TestWhenPlayerCastlesChange_CreatesExactlyOneOwnedCastle(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	zones := applyPlayerCastleChange()

	// Assert
	ownedCount := 0
	for _, mainObject := range zones[0].MainObjects[1:] {
		if mainObject.Owner == "Player1" {
			ownedCount++
		}
	}
	assert.Equal(t, 1, ownedCount, "expected exactly one player-owned extra castle")
}

func TestWhenPlayerCastlesChange_LeavesNeutralZoneUntouched(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	zones := applyPlayerCastleChange()

	// Assert
	assert.Equal(t, 1, test_helpers.NewZoneEditorService().CountZoneCastles(zones[1]),
		"neutral zones must not react to a player option")
}

func TestWhenSpawnZoneLacksSpawnCastle_LeavesItUntouched(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.ZoneConfiguration.PlayerZoneCastles = 2
	zones := []entities.Zone{
		{Name: "Spawn-A", MainObjects: []entities.MainObject{{Type: "City"}}},
	}

	// Act
	newManualReapplyService().ApplyCastleSettingChanges(
		zones, editor_state_dto.CastleSettingChanges{PlayerCastles: true}, configuration)

	// Assert
	assert.Equal(t, []entities.MainObject{{Type: "City"}}, zones[0].MainObjects,
		"a spawn zone without a spawn castle must not be rebuilt")
}

func TestWhenHubCountChanges_RebuildsHubZoneCastles(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.ZoneConfiguration.Advanced.HubZoneCastles = 3
	zones := []entities.Zone{
		{Name: "Hub", MainObjects: []entities.MainObject{{Type: "City"}}},
	}

	// Act
	newManualReapplyService().ApplyCastleSettingChanges(
		zones, editor_state_dto.CastleSettingChanges{Hub: true}, configuration)

	// Assert
	assert.Equal(t, 3, test_helpers.NewZoneEditorService().CountZoneCastles(zones[0]))
}

func TestWhenHubZoneHasLetterSuffix_RebuildsItsCastlesToo(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.ZoneConfiguration.Advanced.HubZoneCastles = 2
	zones := []entities.Zone{
		{Name: "Hub-A", MainObjects: []entities.MainObject{{Type: "City"}}},
	}

	// Act
	newManualReapplyService().ApplyCastleSettingChanges(
		zones, editor_state_dto.CastleSettingChanges{Hub: true}, configuration)

	// Assert
	assert.Equal(t, 2, test_helpers.NewZoneEditorService().CountZoneCastles(zones[0]))
}

func TestWhenNeutralCastlesAreRebuilt_CreatesCastleRoadsToEachExtraCastle(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.ZoneConfiguration.NeutralZoneCastles = 3
	zones := []entities.Zone{makeNeutralZone("G", neutral_zone.QualityMedium, 1)}

	// Act
	newManualReapplyService().ApplyCastleSettingChanges(
		zones, editor_state_dto.CastleSettingChanges{NeutralSimple: true}, configuration)

	// Assert
	castleRoadCount := 0
	for _, road := range zones[0].Roads {
		if road.From.Type == "MainObject" && road.To.Type == "MainObject" {
			castleRoadCount++
		}
	}
	assert.Equal(t, 2, castleRoadCount,
		"expected castle roads from the primary castle to both extra castles")
}
