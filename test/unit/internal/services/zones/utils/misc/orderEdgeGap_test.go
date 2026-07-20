package misc_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones/utils"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenPlayerIsNotAtEnd_OrdersStrongestZoneFirst(t *testing.T) {
	t.Parallel()
	// Arrange
	strongZone := neutral_zone.Plan{Label: "A", Quality: neutral_zone.QualityHigh, CastleCount: gofakeit.Number(0, 4)}
	mediumZone := neutral_zone.Plan{Label: "B", Quality: neutral_zone.QualityMedium, CastleCount: gofakeit.Number(0, 4)}
	weakZone := neutral_zone.Plan{Label: "C", Quality: neutral_zone.QualityLow, CastleCount: gofakeit.Number(0, 4)}
	neutralZones := neutral_zone.Plans{weakZone, strongZone, mediumZone}
	expected := neutral_zone.Plans{strongZone, mediumZone, weakZone}

	// Act
	ordered := utils.OrderEdgeGap(neutralZones, false)

	// Assert
	assert.Equal(t, expected, ordered)
}

func TestWhenPlayerIsAtEnd_OrdersStrongestZoneLast(t *testing.T) {
	t.Parallel()
	// Arrange
	strongZone := neutral_zone.Plan{Label: "A", Quality: neutral_zone.QualityHigh, CastleCount: gofakeit.Number(0, 4)}
	mediumZone := neutral_zone.Plan{Label: "B", Quality: neutral_zone.QualityMedium, CastleCount: gofakeit.Number(0, 4)}
	weakZone := neutral_zone.Plan{Label: "C", Quality: neutral_zone.QualityLow, CastleCount: gofakeit.Number(0, 4)}
	neutralZones := neutral_zone.Plans{weakZone, strongZone, mediumZone}
	expected := neutral_zone.Plans{weakZone, mediumZone, strongZone}

	// Act
	ordered := utils.OrderEdgeGap(neutralZones, true)

	// Assert
	assert.Equal(t, expected, ordered)
}

func TestWhenGapIsEmptyAndPlayerIsAtEnd_ReturnsEmptyPlans(t *testing.T) {
	t.Parallel()
	// Arrange
	neutralZones := neutral_zone.Plans{}

	// Act
	ordered := utils.OrderEdgeGap(neutralZones, true)

	// Assert
	assert.Equal(t, neutral_zone.Plans{}, ordered)
}

func TestWhenGapHasSingleZoneAndPlayerIsAtEnd_ReturnsThatZone(t *testing.T) {
	t.Parallel()
	// Arrange
	zone := neutral_zone.Plan{Label: "A", Quality: neutral_zone.QualityMedium, CastleCount: gofakeit.Number(0, 4)}
	neutralZones := neutral_zone.Plans{zone}

	// Act
	ordered := utils.OrderEdgeGap(neutralZones, true)

	// Assert
	assert.Equal(t, neutral_zone.Plans{zone}, ordered)
}
