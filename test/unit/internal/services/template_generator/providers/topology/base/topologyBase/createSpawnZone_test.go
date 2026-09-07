package topologyBase_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/base"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenSpawnZoneIsCreated_NameCombinesSpawnPrefixWithLabel(t *testing.T) {
	t.Parallel()
	// Arrange
	label := gofakeit.LetterN(3)
	topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())

	// Act
	zone := topologyBase.CreateSpawnZone(newSpawnRequest(label, nil, 1, 0, true, newUnitTuning()))

	// Assert
	assert.Equal(t, "Spawn-"+label, zone.Name)
}

func TestWhenSpawnZoneIsCreated_FirstMainObjectIsSpawnCastleForPlayer(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())
	expected := template_model.MainObject{
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
	zone := topologyBase.CreateSpawnZone(newSpawnRequest("A", nil, 0, 0, true, newUnitTuning()))

	// Assert
	assert.Equal(t, expected, zone.MainObjects[0])
}

func TestWhenOwnedAndUnclaimedCastlesRequested_MainObjectCountIsSpawnPlusOwnedPlusUnclaimed(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())
	tuning := newUnitTuning()
	tuning.PlayerOwnedCastles = 1

	// Act
	zone := topologyBase.CreateSpawnZone(newSpawnRequest("A", nil, 2, 0, true, tuning))

	// Assert
	assert.Len(t, zone.MainObjects, 4)
}

func TestWhenZoneHasNoExtraCastles_RoadsChainConnectionsInsteadOfCastles(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())
	expectedRoads := []template_model.Road{
		{
			From: template_model.TypedRef{Type: "Connection", Args: []string{"Gate-1"}},
			To:   template_model.TypedRef{Type: "Connection", Args: []string{"Gate-2"}},
		},
	}

	// Act
	zone := topologyBase.CreateSpawnZone(
		newSpawnRequest("A", []string{"Gate-1", "Gate-2"}, 0, 0, true, newUnitTuning()))

	// Assert
	assert.Equal(t, expectedRoads, zone.Roads)
}

func TestWhenExtraCastlesArePresent_EveryExtraCastleGetsStoneRoadFromSpawnCastle(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())
	expectedRoads := []template_model.Road{
		{
			Type: "Stone",
			From: template_model.TypedRef{Type: "MainObject", Args: []string{"0"}},
			To:   template_model.TypedRef{Type: "MainObject", Args: []string{"1"}},
		},
		{
			Type: "Stone",
			From: template_model.TypedRef{Type: "MainObject", Args: []string{"0"}},
			To:   template_model.TypedRef{Type: "MainObject", Args: []string{"2"}},
		},
	}

	// Act
	zone := topologyBase.CreateSpawnZone(newSpawnRequest("A", nil, 2, 0, true, newUnitTuning()))

	// Assert
	assert.Equal(t, expectedRoads, zone.Roads)
}

func TestWhenFootholdCountIsPositive_AddsRoadToEveryRemoteFoothold(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())
	expectedRoads := []template_model.Road{
		{
			Type: "Stone",
			From: template_model.TypedRef{Type: "MainObject", Args: []string{"0"}},
			To:   template_model.TypedRef{Type: "MainObject", Args: []string{"1"}},
		},
		{
			From: template_model.TypedRef{Type: "MainObject", Args: []string{"0"}},
			To:   template_model.TypedRef{Type: "MandatoryContent", Args: []string{"name_remote_foothold_1"}},
		},
		{
			From: template_model.TypedRef{Type: "MainObject", Args: []string{"0"}},
			To:   template_model.TypedRef{Type: "MandatoryContent", Args: []string{"name_remote_foothold_2"}},
		},
	}

	// Act
	zone := topologyBase.CreateSpawnZone(newSpawnRequest("A", nil, 1, 2, true, newUnitTuning()))

	// Assert
	assert.Equal(t, expectedRoads, zone.Roads)
}

func TestWhenRoadGenerationIsDisabled_ZoneHasNoRoads(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())

	// Act
	zone := topologyBase.CreateSpawnZone(
		newSpawnRequest("A", []string{"Gate-1"}, 2, 1, false, newUnitTuning()))

	// Assert
	assert.Nil(t, zone.Roads)
}
