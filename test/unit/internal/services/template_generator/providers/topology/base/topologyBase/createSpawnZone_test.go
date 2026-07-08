package topologyBase_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/base"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenSpawnZoneIsCreated_NameCombinesSpawnPrefixWithLabel(t *testing.T) {
	// Arrange
	label := gofakeit.LetterN(3)
	topologyBase := base.NewTopologyBase()

	// Act
	zone := topologyBase.CreateSpawnZone(label, "Player1", nil, 1, false, 1.0, 0, true, newUnitTuning())

	// Assert
	assert.Equal(t, "Spawn-"+label, zone.Name)
}

func TestWhenSpawnZoneIsCreated_FirstMainObjectIsSpawnCastleForPlayer(t *testing.T) {
	// Arrange
	topologyBase := base.NewTopologyBase()
	expected := entities.MainObject{
		Type:                     "Spawn",
		Spawn:                    "Player1",
		RemoveGuardIfHasOwner:    true,
		GuardChance:              1,
		GuardValue:               5000,
		GuardWeeklyIncrement:     0.10,
		BuildingsConstructionSid: "default_buildings_construction",
		Placement:                "Uniform",
		PlacementArgs:            []string{"true", "0.7", "0"},
	}

	// Act
	zone := topologyBase.CreateSpawnZone("A", "Player1", nil, 0, false, 1.0, 0, true, newUnitTuning())

	// Assert
	assert.Equal(t, expected, zone.MainObjects[0])
}

func TestWhenOwnedAndUnclaimedCastlesRequested_MainObjectCountIsSpawnPlusOwnedPlusUnclaimed(t *testing.T) {
	// Arrange
	topologyBase := base.NewTopologyBase()
	tuning := newUnitTuning()
	tuning.PlayerOwnedCastles = 1

	// Act
	zone := topologyBase.CreateSpawnZone("A", "Player1", nil, 2, false, 1.0, 0, true, tuning)

	// Assert
	assert.Len(t, zone.MainObjects, 4)
}

func TestWhenZoneHasNoExtraCastles_RoadsChainConnectionsInsteadOfCastles(t *testing.T) {
	// Arrange
	topologyBase := base.NewTopologyBase()
	expectedRoads := []entities.Road{
		{
			From: entities.TypedRef{Type: "Connection", Args: []string{"Gate-1"}},
			To:   entities.TypedRef{Type: "Connection", Args: []string{"Gate-2"}},
		},
	}

	// Act
	zone := topologyBase.CreateSpawnZone(
		"A", "Player1", []string{"Gate-1", "Gate-2"}, 0, false, 1.0, 0, true, newUnitTuning())

	// Assert
	assert.Equal(t, expectedRoads, zone.Roads)
}

func TestWhenExtraCastlesArePresent_EveryExtraCastleGetsStoneRoadFromSpawnCastle(t *testing.T) {
	// Arrange
	topologyBase := base.NewTopologyBase()
	expectedRoads := []entities.Road{
		{
			Type: "Stone",
			From: entities.TypedRef{Type: "MainObject", Args: []string{"0"}},
			To:   entities.TypedRef{Type: "MainObject", Args: []string{"1"}},
		},
		{
			Type: "Stone",
			From: entities.TypedRef{Type: "MainObject", Args: []string{"0"}},
			To:   entities.TypedRef{Type: "MainObject", Args: []string{"2"}},
		},
	}

	// Act
	zone := topologyBase.CreateSpawnZone("A", "Player1", nil, 2, false, 1.0, 0, true, newUnitTuning())

	// Assert
	assert.Equal(t, expectedRoads, zone.Roads)
}

func TestWhenFootholdCountIsPositive_AddsRoadToEveryRemoteFoothold(t *testing.T) {
	// Arrange
	topologyBase := base.NewTopologyBase()
	expectedRoads := []entities.Road{
		{
			Type: "Stone",
			From: entities.TypedRef{Type: "MainObject", Args: []string{"0"}},
			To:   entities.TypedRef{Type: "MainObject", Args: []string{"1"}},
		},
		{
			From: entities.TypedRef{Type: "MainObject", Args: []string{"0"}},
			To:   entities.TypedRef{Type: "MandatoryContent", Args: []string{"name_remote_foothold_1"}},
		},
		{
			From: entities.TypedRef{Type: "MainObject", Args: []string{"0"}},
			To:   entities.TypedRef{Type: "MandatoryContent", Args: []string{"name_remote_foothold_2"}},
		},
	}

	// Act
	zone := topologyBase.CreateSpawnZone("A", "Player1", nil, 1, false, 1.0, 2, true, newUnitTuning())

	// Assert
	assert.Equal(t, expectedRoads, zone.Roads)
}

func TestWhenRoadGenerationIsDisabled_ZoneHasNoRoads(t *testing.T) {
	// Arrange
	topologyBase := base.NewTopologyBase()

	// Act
	zone := topologyBase.CreateSpawnZone(
		"A", "Player1", []string{"Gate-1"}, 2, false, 1.0, 1, false, newUnitTuning())

	// Assert
	assert.Nil(t, zone.Roads)
}
