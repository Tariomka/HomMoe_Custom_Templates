package misc_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutralZone"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones/utils"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenPlayerIsNotAtEnd_OrdersStrongestZoneFirst(t *testing.T) {
	t.Parallel()
	// Arrange
	strongZone := neutralZone.Plan{Label: "A", Quality: neutralZone.QualityHigh, CastleCount: gofakeit.Number(0, 4)}
	mediumZone := neutralZone.Plan{Label: "B", Quality: neutralZone.QualityMedium, CastleCount: gofakeit.Number(0, 4)}
	weakZone := neutralZone.Plan{Label: "C", Quality: neutralZone.QualityLow, CastleCount: gofakeit.Number(0, 4)}
	neutralZones := neutralZone.Plans{weakZone, strongZone, mediumZone}
	expected := neutralZone.Plans{strongZone, mediumZone, weakZone}

	// Act
	ordered := utils.OrderEdgeGap(neutralZones, false)

	// Assert
	assert.Equal(t, expected, ordered)
}

func TestWhenPlayerIsAtEnd_OrdersStrongestZoneLast(t *testing.T) {
	t.Parallel()
	// Arrange
	strongZone := neutralZone.Plan{Label: "A", Quality: neutralZone.QualityHigh, CastleCount: gofakeit.Number(0, 4)}
	mediumZone := neutralZone.Plan{Label: "B", Quality: neutralZone.QualityMedium, CastleCount: gofakeit.Number(0, 4)}
	weakZone := neutralZone.Plan{Label: "C", Quality: neutralZone.QualityLow, CastleCount: gofakeit.Number(0, 4)}
	neutralZones := neutralZone.Plans{weakZone, strongZone, mediumZone}
	expected := neutralZone.Plans{weakZone, mediumZone, strongZone}

	// Act
	ordered := utils.OrderEdgeGap(neutralZones, true)

	// Assert
	assert.Equal(t, expected, ordered)
}

func TestWhenGapIsEmptyAndPlayerIsAtEnd_ReturnsEmptyPlans(t *testing.T) {
	t.Parallel()
	// Arrange
	neutralZones := neutralZone.Plans{}

	// Act
	ordered := utils.OrderEdgeGap(neutralZones, true)

	// Assert
	assert.Equal(t, neutralZone.Plans{}, ordered)
}

func TestWhenGapHasSingleZoneAndPlayerIsAtEnd_ReturnsThatZone(t *testing.T) {
	t.Parallel()
	// Arrange
	zone := neutralZone.Plan{Label: "A", Quality: neutralZone.QualityMedium, CastleCount: gofakeit.Number(0, 4)}
	neutralZones := neutralZone.Plans{zone}

	// Act
	ordered := utils.OrderEdgeGap(neutralZones, true)

	// Assert
	assert.Equal(t, neutralZone.Plans{zone}, ordered)
}
