package castleFactory_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
	"github.com/stretchr/testify/assert"
)

func TestWhenFactionMatchingIsEnabled_CreatesOwnedMatchFactionCastle(t *testing.T) {
	t.Parallel()
	// Arrange
	factory := zones.NewCastleFactory()

	// Act
	castles := factory.CreatePlayerOwnedCastles(true, "Player1", 1)

	// Assert
	assert.Equal(t, []entities.MainObject{{
		Type:                     "City",
		Owner:                    "Player1",
		BuildingsConstructionSid: "poor_buildings_construction",
		Placement:                "Uniform",
		Faction:                  &entities.TypedRef{Type: "Match", Args: []string{"0"}},
	}}, castles)
}
