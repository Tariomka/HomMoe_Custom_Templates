package topologyBase_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/base"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenOwnedCastleCountIsZero_NoOwnedCastlesAreCreated(t *testing.T) {
	// Arrange
	topologyBase := base.NewTopologyBase()

	// Act
	castles := topologyBase.CreatePlayerOwnedCastles(false, "Player1", 0)

	// Assert
	assert.Empty(t, castles)
}

func TestWhenFactionsMatchPlayer_OwnedCastleIsPoorCityWithMatchFaction(t *testing.T) {
	// Arrange
	topologyBase := base.NewTopologyBase()
	expectedCastles := []entities.MainObject{
		{
			Type:                     "City",
			Owner:                    "Player1",
			BuildingsConstructionSid: "poor_buildings_construction",
			Placement:                "Uniform",
			Faction:                  &entities.TypedRef{Type: "Match", Args: []string{"0"}},
		},
	}

	// Act
	castles := topologyBase.CreatePlayerOwnedCastles(true, "Player1", 1)

	// Assert
	assert.Equal(t, expectedCastles, castles)
}

func TestWhenFactionsDoNotMatchPlayer_OwnedCastleFactionIsRandom(t *testing.T) {
	// Arrange
	topologyBase := base.NewTopologyBase()
	expectedCastles := []entities.MainObject{
		{
			Type:                     "City",
			Owner:                    "Player1",
			BuildingsConstructionSid: "poor_buildings_construction",
			Placement:                "Uniform",
			Faction:                  &entities.TypedRef{Type: "Random"},
		},
	}

	// Act
	castles := topologyBase.CreatePlayerOwnedCastles(false, "Player1", 1)

	// Assert
	assert.Equal(t, expectedCastles, castles)
}

func TestWhenMultipleOwnedCastlesAreRequested_CastleCountMatchesRequest(t *testing.T) {
	// Arrange
	topologyBase := base.NewTopologyBase()
	castleCount := gofakeit.Number(2, 5)

	// Act
	castles := topologyBase.CreatePlayerOwnedCastles(false, "Player1", castleCount)

	// Assert
	assert.Len(t, castles, castleCount)
}
