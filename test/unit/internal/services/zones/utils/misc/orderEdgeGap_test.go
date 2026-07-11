package misc_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones/utils"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenPlayerIsNotAtEnd_OrdersStrongestZoneFirst(t *testing.T) {
	t.Parallel()
	// Arrange
	strongZone := models.NeutralZonePlan{Label: "A", Quality: models.QualityHigh, CastleCount: gofakeit.Number(0, 4)}
	mediumZone := models.NeutralZonePlan{Label: "B", Quality: models.QualityMedium, CastleCount: gofakeit.Number(0, 4)}
	weakZone := models.NeutralZonePlan{Label: "C", Quality: models.QualityLow, CastleCount: gofakeit.Number(0, 4)}
	neutralZones := models.NeutralZonePlans{weakZone, strongZone, mediumZone}
	expected := models.NeutralZonePlans{strongZone, mediumZone, weakZone}

	// Act
	ordered := utils.OrderEdgeGap(neutralZones, false)

	// Assert
	assert.Equal(t, expected, ordered)
}

func TestWhenPlayerIsAtEnd_OrdersStrongestZoneLast(t *testing.T) {
	t.Parallel()
	// Arrange
	strongZone := models.NeutralZonePlan{Label: "A", Quality: models.QualityHigh, CastleCount: gofakeit.Number(0, 4)}
	mediumZone := models.NeutralZonePlan{Label: "B", Quality: models.QualityMedium, CastleCount: gofakeit.Number(0, 4)}
	weakZone := models.NeutralZonePlan{Label: "C", Quality: models.QualityLow, CastleCount: gofakeit.Number(0, 4)}
	neutralZones := models.NeutralZonePlans{weakZone, strongZone, mediumZone}
	expected := models.NeutralZonePlans{weakZone, mediumZone, strongZone}

	// Act
	ordered := utils.OrderEdgeGap(neutralZones, true)

	// Assert
	assert.Equal(t, expected, ordered)
}

func TestWhenGapIsEmptyAndPlayerIsAtEnd_ReturnsEmptyPlans(t *testing.T) {
	t.Parallel()
	// Arrange
	neutralZones := models.NeutralZonePlans{}

	// Act
	ordered := utils.OrderEdgeGap(neutralZones, true)

	// Assert
	assert.Equal(t, models.NeutralZonePlans{}, ordered)
}

func TestWhenGapHasSingleZoneAndPlayerIsAtEnd_ReturnsThatZone(t *testing.T) {
	t.Parallel()
	// Arrange
	zone := models.NeutralZonePlan{Label: "A", Quality: models.QualityMedium, CastleCount: gofakeit.Number(0, 4)}
	neutralZones := models.NeutralZonePlans{zone}

	// Act
	ordered := utils.OrderEdgeGap(neutralZones, true)

	// Assert
	assert.Equal(t, models.NeutralZonePlans{zone}, ordered)
}
