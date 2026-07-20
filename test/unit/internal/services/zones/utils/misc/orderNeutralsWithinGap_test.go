package misc_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones/utils"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenGapIsEmpty_ReturnsEmptyPlans(t *testing.T) {
	t.Parallel()
	// Arrange
	neutralZones := neutral_zone.Plans{}

	// Act
	ordered := utils.OrderNeutralsWithinGap(neutralZones)

	// Assert
	assert.Equal(t, neutral_zone.Plans{}, ordered)
}

func TestWhenGapHasSingleZone_ReturnsThatZone(t *testing.T) {
	t.Parallel()
	// Arrange
	zone := neutral_zone.Plan{Label: "A", Quality: neutral_zone.QualityMedium, CastleCount: gofakeit.Number(0, 4)}
	neutralZones := neutral_zone.Plans{zone}

	// Act
	ordered := utils.OrderNeutralsWithinGap(neutralZones)

	// Assert
	assert.Equal(t, neutral_zone.Plans{zone}, ordered)
}

func TestWhenGapHasThreeZones_PlacesStrongestFirstAndSecondStrongestLast(t *testing.T) {
	t.Parallel()
	// Arrange
	strongZone := neutral_zone.Plan{Label: "A", Quality: neutral_zone.QualityHigh, CastleCount: gofakeit.Number(0, 4)}
	mediumZone := neutral_zone.Plan{Label: "B", Quality: neutral_zone.QualityMedium, CastleCount: gofakeit.Number(0, 4)}
	weakZone := neutral_zone.Plan{Label: "C", Quality: neutral_zone.QualityLow, CastleCount: gofakeit.Number(0, 4)}
	neutralZones := neutral_zone.Plans{weakZone, strongZone, mediumZone}
	expected := neutral_zone.Plans{strongZone, weakZone, mediumZone}

	// Act
	ordered := utils.OrderNeutralsWithinGap(neutralZones)

	// Assert
	assert.Equal(t, expected, ordered)
}

func TestWhenGapHasFourZones_AlternatesEndsInwards(t *testing.T) {
	t.Parallel()
	// Arrange
	firstZone := neutral_zone.Plan{Label: "A", Quality: neutral_zone.QualityHigh, CastleCount: 4}
	secondZone := neutral_zone.Plan{Label: "B", Quality: neutral_zone.QualityHigh, CastleCount: 0}
	thirdZone := neutral_zone.Plan{Label: "C", Quality: neutral_zone.QualityLow, CastleCount: 4}
	fourthZone := neutral_zone.Plan{Label: "D", Quality: neutral_zone.QualityLow, CastleCount: 0}
	neutralZones := neutral_zone.Plans{fourthZone, secondZone, firstZone, thirdZone}
	expected := neutral_zone.Plans{firstZone, thirdZone, fourthZone, secondZone}

	// Act
	ordered := utils.OrderNeutralsWithinGap(neutralZones)

	// Assert
	assert.Equal(t, expected, ordered)
}
