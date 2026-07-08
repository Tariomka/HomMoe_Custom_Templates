package topologyBase_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/base"
	"github.com/stretchr/testify/assert"
)

func TestWhenNeutralCastleCountIsZero_NoNeutralCastlesAreCreated(t *testing.T) {
	// Arrange
	profile := models.NewNeutralZoneProfile(models.QualityMedium)

	// Act
	castles := base.CreateNeutralZoneCastles(profile, newUnitTuning(), 0, false)

	// Assert
	assert.Empty(t, castles)
}

func TestWhenZoneIsHoldCity_PrimaryCastleIsUltraRichCenteredHoldCityCastle(t *testing.T) {
	// Arrange
	profile := models.NewNeutralZoneProfile(models.QualityHigh)
	expectedCastles := []entities.MainObject{
		{
			Type:                     "City",
			GuardChance:              1,
			GuardValue:               20000,
			GuardWeeklyIncrement:     0.10,
			BuildingsConstructionSid: "ultra_rich_buildings_construction",
			Faction:                  &entities.TypedRef{Type: "FromList"},
			Placement:                "Center",
			HoldCityWinCon:           true,
		},
	}

	// Act
	castles := base.CreateNeutralZoneCastles(profile, newUnitTuning(), 1, true)

	// Assert
	assert.Equal(t, expectedCastles, castles)
}

func TestWhenHoldCityProfileGuardExceedsFloor_ProfileGuardValueIsUsed(t *testing.T) {
	// Arrange
	profile := models.NeutralZoneProfile{PrimaryCityGuardValue: 30000}

	// Act
	castles := base.CreateNeutralZoneCastles(profile, newUnitTuning(), 1, true)

	// Assert
	assert.Equal(t, 30000, castles[0].GuardValue)
}

func TestWhenZoneIsNotHoldCity_PrimaryCastleUsesProfileQualityAndGuard(t *testing.T) {
	// Arrange
	profile := models.NewNeutralZoneProfile(models.QualityLow)
	expectedCastles := []entities.MainObject{
		{
			Type:                     "City",
			GuardChance:              1,
			GuardValue:               4000,
			GuardWeeklyIncrement:     0.10,
			BuildingsConstructionSid: "poor_buildings_construction",
			Faction:                  &entities.TypedRef{Type: "FromList"},
			Placement:                "Uniform",
			PlacementArgs:            []string{"true", "0.8", "2"},
		},
	}

	// Act
	castles := base.CreateNeutralZoneCastles(profile, newUnitTuning(), 1, false)

	// Assert
	assert.Equal(t, expectedCastles, castles)
}

func TestWhenMultipleNeutralCastlesAreRequested_ExtraCastlesUseExtraProfileValues(t *testing.T) {
	// Arrange
	profile := models.NewNeutralZoneProfile(models.QualityMedium)
	expectedExtraCastle := entities.MainObject{
		Type:                     "City",
		GuardChance:              1,
		GuardValue:               4000,
		GuardWeeklyIncrement:     0.10,
		BuildingsConstructionSid: "poor_buildings_construction",
		Faction:                  &entities.TypedRef{Type: "FromList"},
		Placement:                "Uniform",
		PlacementArgs:            []string{"false", "-0.8", "3"},
	}

	// Act
	castles := base.CreateNeutralZoneCastles(profile, newUnitTuning(), 3, false)

	// Assert
	assert.Equal(t, []entities.MainObject{expectedExtraCastle, expectedExtraCastle}, castles[1:])
}

func TestWhenBorderGuardMultiplierIsDoubled_NeutralCastleGuardsAreScaled(t *testing.T) {
	// Arrange
	profile := models.NewNeutralZoneProfile(models.QualityLow)
	tuning := newUnitTuning()
	tuning.BorderGuardStrengthMultiplier = 2.0

	// Act
	castles := base.CreateNeutralZoneCastles(profile, tuning, 1, false)

	// Assert
	assert.Equal(t, 8000, castles[0].GuardValue)
}
