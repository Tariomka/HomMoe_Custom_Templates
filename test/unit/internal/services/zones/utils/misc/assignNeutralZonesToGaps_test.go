package misc_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones/utils"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenCapacitiesAreEmpty_ReturnsNoGaps(t *testing.T) {
	// Arrange
	neutralZones := models.NeutralZonePlans{
		{Label: "A", Quality: models.QualityHigh, CastleCount: gofakeit.Number(0, 4)},
	}

	// Act
	gaps := utils.AssignNeutralZonesToGaps(neutralZones, []int{}, false)

	// Assert
	assert.Empty(t, gaps)
}

func TestWhenZonesFitExactly_AssignsStrongestZoneToLowestIndexedGap(t *testing.T) {
	// Arrange
	strongZone := models.NeutralZonePlan{Label: "A", Quality: models.QualityHigh, CastleCount: gofakeit.Number(0, 4)}
	weakZone := models.NeutralZonePlan{Label: "B", Quality: models.QualityMedium, CastleCount: gofakeit.Number(0, 4)}
	neutralZones := models.NeutralZonePlans{weakZone, strongZone}
	expected := []models.NeutralZonePlans{{strongZone}, {weakZone}}

	// Act
	gaps := utils.AssignNeutralZonesToGaps(neutralZones, []int{1, 1}, false)

	// Assert
	assert.Equal(t, expected, gaps)
}

func TestWhenTotalCapacityIsExceeded_DropsWeakestZones(t *testing.T) {
	// Arrange
	strongZone := models.NeutralZonePlan{Label: "A", Quality: models.QualityHigh, CastleCount: gofakeit.Number(0, 4)}
	weakZone := models.NeutralZonePlan{Label: "B", Quality: models.QualityLow, CastleCount: gofakeit.Number(0, 4)}
	neutralZones := models.NeutralZonePlans{weakZone, strongZone}
	expected := []models.NeutralZonePlans{{strongZone}}

	// Act
	gaps := utils.AssignNeutralZonesToGaps(neutralZones, []int{1}, false)

	// Assert
	assert.Equal(t, expected, gaps)
}

func TestWhenPreferInteriorIsTrue_AssignsFirstZoneToInteriorGap(t *testing.T) {
	// Arrange
	zone := models.NeutralZonePlan{Label: "A", Quality: models.QualityHigh, CastleCount: gofakeit.Number(0, 4)}
	neutralZones := models.NeutralZonePlans{zone}
	expected := []models.NeutralZonePlans{nil, {zone}, nil}

	// Act
	gaps := utils.AssignNeutralZonesToGaps(neutralZones, []int{1, 1, 1}, true)

	// Assert
	assert.Equal(t, expected, gaps)
}

func TestWhenInteriorGapsAreFull_FallsBackToEdgeGaps(t *testing.T) {
	// Arrange
	strongZone := models.NeutralZonePlan{Label: "A", Quality: models.QualityHigh, CastleCount: gofakeit.Number(0, 4)}
	mediumZone := models.NeutralZonePlan{Label: "B", Quality: models.QualityMedium, CastleCount: gofakeit.Number(0, 4)}
	weakZone := models.NeutralZonePlan{Label: "C", Quality: models.QualityLow, CastleCount: gofakeit.Number(0, 4)}
	neutralZones := models.NeutralZonePlans{weakZone, mediumZone, strongZone}
	expected := []models.NeutralZonePlans{{mediumZone}, {strongZone}, {weakZone}}

	// Act
	gaps := utils.AssignNeutralZonesToGaps(neutralZones, []int{1, 1, 1}, true)

	// Assert
	assert.Equal(t, expected, gaps)
}

func TestWhenOnlyTwoGapsExistWithPreferInterior_UsesEdgeGap(t *testing.T) {
	// Arrange
	zone := models.NeutralZonePlan{Label: "A", Quality: models.QualityMedium, CastleCount: gofakeit.Number(0, 4)}
	neutralZones := models.NeutralZonePlans{zone}
	expected := []models.NeutralZonePlans{{zone}, nil}

	// Act
	gaps := utils.AssignNeutralZonesToGaps(neutralZones, []int{1, 1}, true)

	// Assert
	assert.Equal(t, expected, gaps)
}

func TestWhenGapLoadsAreTied_PrefersGapWithFewerZones(t *testing.T) {
	// Arrange - one medium zone (score 2.0) balances two low zones (1.0 each),
	// so the fourth zone sees equal loads but fewer zones in the first gap.
	mediumZone := models.NeutralZonePlan{Label: "A", Quality: models.QualityMedium, CastleCount: 0}
	firstLowZone := models.NeutralZonePlan{Label: "B", Quality: models.QualityLow, CastleCount: 0}
	secondLowZone := models.NeutralZonePlan{Label: "C", Quality: models.QualityLow, CastleCount: 0}
	tieBreakerZone := models.NeutralZonePlan{Label: "D", Quality: models.QualityLow, CastleCount: 0}
	neutralZones := models.NeutralZonePlans{mediumZone, firstLowZone, secondLowZone, tieBreakerZone}
	expected := []models.NeutralZonePlans{{mediumZone, tieBreakerZone}, {firstLowZone, secondLowZone}}

	// Act
	gaps := utils.AssignNeutralZonesToGaps(neutralZones, []int{3, 3}, false)

	// Assert
	assert.Equal(t, expected, gaps)
}

func TestWhenNoZonesAreGiven_ReturnsEmptyGapPerCapacity(t *testing.T) {
	// Arrange
	expected := []models.NeutralZonePlans{nil, nil}

	// Act
	gaps := utils.AssignNeutralZonesToGaps(models.NeutralZonePlans{}, []int{1, 1}, false)

	// Assert
	assert.Equal(t, expected, gaps)
}
