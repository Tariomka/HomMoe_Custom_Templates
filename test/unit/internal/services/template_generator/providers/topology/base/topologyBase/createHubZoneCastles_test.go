package topologyBase_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/base"
	"github.com/stretchr/testify/assert"
)

func TestWhenHubCastleCountIsZero_NoHubCastlesAreCreatedEvenForHoldCity(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase()

	// Act
	castles := topologyBase.CreateHubZoneCastles(newUnitTuning(), 0, true)

	// Assert
	assert.Empty(t, castles)
}

func TestWhenHubIsHoldCity_FirstCastleIsUltraRichCenteredHoldCityCastle(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase()
	expectedCastles := []entities.MainObject{
		{
			Type:                     "City",
			GuardChance:              1,
			GuardValue:               25000,
			GuardWeeklyIncrement:     0.10,
			BuildingsConstructionSid: "ultra_rich_buildings_construction",
			Faction:                  &entities.TypedRef{Type: "FromList"},
			Placement:                "Center",
			HoldCityWinCon:           true,
		},
	}

	// Act
	castles := topologyBase.CreateHubZoneCastles(newUnitTuning(), 1, true)

	// Assert
	assert.Equal(t, expectedCastles, castles)
}

func TestWhenHubIsNotHoldCity_CastleIsRichWithHalfGuardChance(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase()
	expectedCastles := []entities.MainObject{
		{
			Type:                     "City",
			GuardChance:              0.5,
			GuardValue:               16000,
			GuardWeeklyIncrement:     0.10,
			BuildingsConstructionSid: "rich_buildings_construction",
			Faction:                  &entities.TypedRef{Type: "FromList"},
			Placement:                "Uniform",
			PlacementArgs:            []string{"true", "0.8", "2"},
		},
	}

	// Act
	castles := topologyBase.CreateHubZoneCastles(newUnitTuning(), 1, false)

	// Assert
	assert.Equal(t, expectedCastles, castles)
}

func TestWhenHoldCityHubHasMultipleCastles_ExtraCastlesAreRegularRichCastles(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase()
	expectedExtraCastle := entities.MainObject{
		Type:                     "City",
		GuardChance:              0.5,
		GuardValue:               16000,
		GuardWeeklyIncrement:     0.10,
		BuildingsConstructionSid: "rich_buildings_construction",
		Faction:                  &entities.TypedRef{Type: "FromList"},
		Placement:                "Uniform",
		PlacementArgs:            []string{"true", "0.8", "2"},
	}

	// Act
	castles := topologyBase.CreateHubZoneCastles(newUnitTuning(), 3, true)

	// Assert
	assert.Equal(t, []entities.MainObject{expectedExtraCastle, expectedExtraCastle}, castles[1:])
}

func TestWhenNeutralGuardMultiplierIsDoubled_HubCastleGuardIsScaled(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase()
	tuning := newUnitTuning()
	tuning.NeutralStackStrengthMultiplier = 2.0

	// Act
	castles := topologyBase.CreateHubZoneCastles(tuning, 1, false)

	// Assert
	assert.Equal(t, 32000, castles[0].GuardValue)
}
