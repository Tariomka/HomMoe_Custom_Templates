package misc_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutralZone"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones/utils"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenGapIsEmpty_ReturnsEmptyPlans(t *testing.T) {
	t.Parallel()
	// Arrange
	neutralZones := neutralZone.Plans{}

	// Act
	ordered := utils.OrderNeutralsWithinGap(neutralZones)

	// Assert
	assert.Equal(t, neutralZone.Plans{}, ordered)
}

func TestWhenGapHasSingleZone_ReturnsThatZone(t *testing.T) {
	t.Parallel()
	// Arrange
	zone := neutralZone.Plan{Label: "A", Quality: neutralZone.QualityMedium, CastleCount: gofakeit.Number(0, 4)}
	neutralZones := neutralZone.Plans{zone}

	// Act
	ordered := utils.OrderNeutralsWithinGap(neutralZones)

	// Assert
	assert.Equal(t, neutralZone.Plans{zone}, ordered)
}

func TestWhenGapHasThreeZones_PlacesStrongestFirstAndSecondStrongestLast(t *testing.T) {
	t.Parallel()
	// Arrange
	strongZone := neutralZone.Plan{Label: "A", Quality: neutralZone.QualityHigh, CastleCount: gofakeit.Number(0, 4)}
	mediumZone := neutralZone.Plan{Label: "B", Quality: neutralZone.QualityMedium, CastleCount: gofakeit.Number(0, 4)}
	weakZone := neutralZone.Plan{Label: "C", Quality: neutralZone.QualityLow, CastleCount: gofakeit.Number(0, 4)}
	neutralZones := neutralZone.Plans{weakZone, strongZone, mediumZone}
	expected := neutralZone.Plans{strongZone, weakZone, mediumZone}

	// Act
	ordered := utils.OrderNeutralsWithinGap(neutralZones)

	// Assert
	assert.Equal(t, expected, ordered)
}

func TestWhenGapHasFourZones_AlternatesEndsInwards(t *testing.T) {
	t.Parallel()
	// Arrange
	firstZone := neutralZone.Plan{Label: "A", Quality: neutralZone.QualityHigh, CastleCount: 4}
	secondZone := neutralZone.Plan{Label: "B", Quality: neutralZone.QualityHigh, CastleCount: 0}
	thirdZone := neutralZone.Plan{Label: "C", Quality: neutralZone.QualityLow, CastleCount: 4}
	fourthZone := neutralZone.Plan{Label: "D", Quality: neutralZone.QualityLow, CastleCount: 0}
	neutralZones := neutralZone.Plans{fourthZone, secondZone, firstZone, thirdZone}
	expected := neutralZone.Plans{firstZone, thirdZone, fourthZone, secondZone}

	// Act
	ordered := utils.OrderNeutralsWithinGap(neutralZones)

	// Assert
	assert.Equal(t, expected, ordered)
}
