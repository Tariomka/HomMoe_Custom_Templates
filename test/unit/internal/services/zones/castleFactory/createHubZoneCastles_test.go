package castleFactory_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
	"github.com/stretchr/testify/assert"
)

func TestWhenHubIsNotHoldCity_CreatesRichHalfChanceCastle(t *testing.T) {
	t.Parallel()
	// Arrange
	factory := zones.NewCastleFactory()

	// Act
	castles := factory.CreateHubZoneCastles(newUnitTuning(), 1, false)

	// Assert
	assert.Equal(t, []template_model.MainObject{{
		Type:                     "City",
		GuardChance:              0.5,
		GuardValue:               16000,
		GuardWeeklyIncrement:     0.10,
		BuildingsConstructionSid: "rich_buildings_construction",
		Faction:                  &template_model.TypedRef{Type: "FromList"},
		Placement:                "Uniform",
		PlacementArgs:            []string{"true", "0.8", "2"},
	}}, castles)
}
