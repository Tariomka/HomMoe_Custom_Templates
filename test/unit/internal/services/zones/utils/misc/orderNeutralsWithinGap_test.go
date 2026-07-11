package misc_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones/utils"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenGapIsEmpty_ReturnsEmptyPlans(t *testing.T) {
	t.Parallel()
	// Arrange
	neutralZones := models.NeutralZonePlans{}

	// Act
	ordered := utils.OrderNeutralsWithinGap(neutralZones)

	// Assert
	assert.Equal(t, models.NeutralZonePlans{}, ordered)
}

func TestWhenGapHasSingleZone_ReturnsThatZone(t *testing.T) {
	t.Parallel()
	// Arrange
	zone := models.NeutralZonePlan{Label: "A", Quality: models.QualityMedium, CastleCount: gofakeit.Number(0, 4)}
	neutralZones := models.NeutralZonePlans{zone}

	// Act
	ordered := utils.OrderNeutralsWithinGap(neutralZones)

	// Assert
	assert.Equal(t, models.NeutralZonePlans{zone}, ordered)
}

func TestWhenGapHasThreeZones_PlacesStrongestFirstAndSecondStrongestLast(t *testing.T) {
	t.Parallel()
	// Arrange
	strongZone := models.NeutralZonePlan{Label: "A", Quality: models.QualityHigh, CastleCount: gofakeit.Number(0, 4)}
	mediumZone := models.NeutralZonePlan{Label: "B", Quality: models.QualityMedium, CastleCount: gofakeit.Number(0, 4)}
	weakZone := models.NeutralZonePlan{Label: "C", Quality: models.QualityLow, CastleCount: gofakeit.Number(0, 4)}
	neutralZones := models.NeutralZonePlans{weakZone, strongZone, mediumZone}
	expected := models.NeutralZonePlans{strongZone, weakZone, mediumZone}

	// Act
	ordered := utils.OrderNeutralsWithinGap(neutralZones)

	// Assert
	assert.Equal(t, expected, ordered)
}

func TestWhenGapHasFourZones_AlternatesEndsInwards(t *testing.T) {
	t.Parallel()
	// Arrange
	firstZone := models.NeutralZonePlan{Label: "A", Quality: models.QualityHigh, CastleCount: 4}
	secondZone := models.NeutralZonePlan{Label: "B", Quality: models.QualityHigh, CastleCount: 0}
	thirdZone := models.NeutralZonePlan{Label: "C", Quality: models.QualityLow, CastleCount: 4}
	fourthZone := models.NeutralZonePlan{Label: "D", Quality: models.QualityLow, CastleCount: 0}
	neutralZones := models.NeutralZonePlans{fourthZone, secondZone, firstZone, thirdZone}
	expected := models.NeutralZonePlans{firstZone, thirdZone, fourthZone, secondZone}

	// Act
	ordered := utils.OrderNeutralsWithinGap(neutralZones)

	// Assert
	assert.Equal(t, expected, ordered)
}
