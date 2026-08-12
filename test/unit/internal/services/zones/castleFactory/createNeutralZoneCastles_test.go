package castleFactory_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_zones"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
	"github.com/stretchr/testify/assert"
)

func TestWhenNeutralZoneIsHoldCity_CreatesCenteredHoldCityCastle(t *testing.T) {
	t.Parallel()
	// Arrange
	factory := zones.NewCastleFactory()
	profile := common_zones.GetNeutralZoneProfile(neutral_zone.QualityHigh)

	// Act
	castles := factory.CreateNeutralZoneCastles(profile, newUnitTuning(), 1, true)

	// Assert
	assert.Equal(t, []entities.MainObject{{
		Type:                     "City",
		GuardChance:              1,
		GuardValue:               20000,
		GuardWeeklyIncrement:     0.10,
		BuildingsConstructionSid: "ultra_rich_buildings_construction",
		Faction:                  &entities.TypedRef{Type: "FromList"},
		Placement:                "Center",
		HoldCityWinCon:           true,
	}}, castles)
}
