package castleFactory_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
	"github.com/stretchr/testify/assert"
)

func TestWhenFactionMatchingIsDisabled_CreatesGuardedRandomFactionCastle(t *testing.T) {
	t.Parallel()
	// Arrange
	factory := zones.NewCastleFactory()

	// Act
	castles := factory.CreatePlayerUnclaimedCastles(false, 5000, 1)

	// Assert
	assert.Equal(t, []entities.MainObject{{
		Type:                     "City",
		GuardChance:              1,
		GuardValue:               5000,
		GuardWeeklyIncrement:     0.15,
		BuildingsConstructionSid: "medium_buildings_construction",
		Faction:                  &entities.TypedRef{Type: "Random"},
		Placement:                "Uniform",
		PlacementArgs:            []string{"false", "-0.8", "3"},
	}}, castles)
}
