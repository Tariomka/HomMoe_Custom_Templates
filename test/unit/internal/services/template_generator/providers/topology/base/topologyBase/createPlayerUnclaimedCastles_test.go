package topologyBase_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/base"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenUnclaimedCastleCountIsZero_NoUnclaimedCastlesAreCreated(t *testing.T) {
	// Arrange
	topologyBase := base.NewTopologyBase()

	// Act
	castles := topologyBase.CreatePlayerUnclaimedCastles(false, 5000, 0)

	// Assert
	assert.Empty(t, castles)
}

func TestWhenFactionsMatchPlayer_UnclaimedCastleIsGuardedMediumCityWithMatchFaction(t *testing.T) {
	// Arrange
	topologyBase := base.NewTopologyBase()
	expectedCastles := []entities.MainObject{
		{
			Type:                     "City",
			GuardChance:              1,
			GuardValue:               5000,
			GuardWeeklyIncrement:     0.15,
			BuildingsConstructionSid: "medium_buildings_construction",
			Placement:                "Uniform",
			PlacementArgs:            []string{"false", "-0.8", "3"},
			Faction:                  &entities.TypedRef{Type: "Match", Args: []string{"0"}},
		},
	}

	// Act
	castles := topologyBase.CreatePlayerUnclaimedCastles(true, 5000, 1)

	// Assert
	assert.Equal(t, expectedCastles, castles)
}

func TestWhenFactionsDoNotMatchPlayer_UnclaimedCastleFactionIsRandom(t *testing.T) {
	// Arrange
	topologyBase := base.NewTopologyBase()
	expectedCastles := []entities.MainObject{
		{
			Type:                     "City",
			GuardChance:              1,
			GuardValue:               5000,
			GuardWeeklyIncrement:     0.15,
			BuildingsConstructionSid: "medium_buildings_construction",
			Placement:                "Uniform",
			PlacementArgs:            []string{"false", "-0.8", "3"},
			Faction:                  &entities.TypedRef{Type: "Random"},
		},
	}

	// Act
	castles := topologyBase.CreatePlayerUnclaimedCastles(false, 5000, 1)

	// Assert
	assert.Equal(t, expectedCastles, castles)
}

func TestWhenGuardValueIsProvided_UnclaimedCastleCarriesItVerbatim(t *testing.T) {
	// Arrange
	topologyBase := base.NewTopologyBase()
	guardValue := gofakeit.Number(1000, 50000)

	// Act
	castles := topologyBase.CreatePlayerUnclaimedCastles(false, guardValue, 1)

	// Assert
	assert.Equal(t, guardValue, castles[0].GuardValue)
}

func TestWhenMultipleUnclaimedCastlesAreRequested_CastleCountMatchesRequest(t *testing.T) {
	// Arrange
	topologyBase := base.NewTopologyBase()
	castleCount := gofakeit.Number(2, 5)

	// Act
	castles := topologyBase.CreatePlayerUnclaimedCastles(false, 5000, castleCount)

	// Assert
	assert.Len(t, castles, castleCount)
}
