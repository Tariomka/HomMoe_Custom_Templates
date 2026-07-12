package misc_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutralZone"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones/utils"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenCapacitiesAreEmpty_ReturnsNoGaps(t *testing.T) {
	t.Parallel()
	// Arrange
	neutralZones := neutralZone.Plans{
		{Label: "A", Quality: neutralZone.QualityHigh, CastleCount: gofakeit.Number(0, 4)},
	}

	// Act
	gaps := utils.AssignNeutralZonesToGaps(neutralZones, []int{}, false)

	// Assert
	assert.Empty(t, gaps)
}

func TestWhenZonesFitExactly_AssignsStrongestZoneToLowestIndexedGap(t *testing.T) {
	t.Parallel()
	// Arrange
	strongZone := neutralZone.Plan{Label: "A", Quality: neutralZone.QualityHigh, CastleCount: gofakeit.Number(0, 4)}
	weakZone := neutralZone.Plan{Label: "B", Quality: neutralZone.QualityMedium, CastleCount: gofakeit.Number(0, 4)}
	neutralZones := neutralZone.Plans{weakZone, strongZone}
	expected := []neutralZone.Plans{{strongZone}, {weakZone}}

	// Act
	gaps := utils.AssignNeutralZonesToGaps(neutralZones, []int{1, 1}, false)

	// Assert
	assert.Equal(t, expected, gaps)
}

func TestWhenTotalCapacityIsExceeded_DropsWeakestZones(t *testing.T) {
	t.Parallel()
	// Arrange
	strongZone := neutralZone.Plan{Label: "A", Quality: neutralZone.QualityHigh, CastleCount: gofakeit.Number(0, 4)}
	weakZone := neutralZone.Plan{Label: "B", Quality: neutralZone.QualityLow, CastleCount: gofakeit.Number(0, 4)}
	neutralZones := neutralZone.Plans{weakZone, strongZone}
	expected := []neutralZone.Plans{{strongZone}}

	// Act
	gaps := utils.AssignNeutralZonesToGaps(neutralZones, []int{1}, false)

	// Assert
	assert.Equal(t, expected, gaps)
}

func TestWhenPreferInteriorIsTrue_AssignsFirstZoneToInteriorGap(t *testing.T) {
	t.Parallel()
	// Arrange
	zone := neutralZone.Plan{Label: "A", Quality: neutralZone.QualityHigh, CastleCount: gofakeit.Number(0, 4)}
	neutralZones := neutralZone.Plans{zone}
	expected := []neutralZone.Plans{nil, {zone}, nil}

	// Act
	gaps := utils.AssignNeutralZonesToGaps(neutralZones, []int{1, 1, 1}, true)

	// Assert
	assert.Equal(t, expected, gaps)
}

func TestWhenInteriorGapsAreFull_FallsBackToEdgeGaps(t *testing.T) {
	t.Parallel()
	// Arrange
	strongZone := neutralZone.Plan{Label: "A", Quality: neutralZone.QualityHigh, CastleCount: gofakeit.Number(0, 4)}
	mediumZone := neutralZone.Plan{Label: "B", Quality: neutralZone.QualityMedium, CastleCount: gofakeit.Number(0, 4)}
	weakZone := neutralZone.Plan{Label: "C", Quality: neutralZone.QualityLow, CastleCount: gofakeit.Number(0, 4)}
	neutralZones := neutralZone.Plans{weakZone, mediumZone, strongZone}
	expected := []neutralZone.Plans{{mediumZone}, {strongZone}, {weakZone}}

	// Act
	gaps := utils.AssignNeutralZonesToGaps(neutralZones, []int{1, 1, 1}, true)

	// Assert
	assert.Equal(t, expected, gaps)
}

func TestWhenOnlyTwoGapsExistWithPreferInterior_UsesEdgeGap(t *testing.T) {
	t.Parallel()
	// Arrange
	zone := neutralZone.Plan{Label: "A", Quality: neutralZone.QualityMedium, CastleCount: gofakeit.Number(0, 4)}
	neutralZones := neutralZone.Plans{zone}
	expected := []neutralZone.Plans{{zone}, nil}

	// Act
	gaps := utils.AssignNeutralZonesToGaps(neutralZones, []int{1, 1}, true)

	// Assert
	assert.Equal(t, expected, gaps)
}

func TestWhenGapLoadsAreTied_PrefersGapWithFewerZones(t *testing.T) {
	t.Parallel()
	// Arrange - one medium zone (score 2.0) balances two low zones (1.0 each),
	// so the fourth zone sees equal loads but fewer zones in the first gap.
	mediumZone := neutralZone.Plan{Label: "A", Quality: neutralZone.QualityMedium, CastleCount: 0}
	firstLowZone := neutralZone.Plan{Label: "B", Quality: neutralZone.QualityLow, CastleCount: 0}
	secondLowZone := neutralZone.Plan{Label: "C", Quality: neutralZone.QualityLow, CastleCount: 0}
	tieBreakerZone := neutralZone.Plan{Label: "D", Quality: neutralZone.QualityLow, CastleCount: 0}
	neutralZones := neutralZone.Plans{mediumZone, firstLowZone, secondLowZone, tieBreakerZone}
	expected := []neutralZone.Plans{{mediumZone, tieBreakerZone}, {firstLowZone, secondLowZone}}

	// Act
	gaps := utils.AssignNeutralZonesToGaps(neutralZones, []int{3, 3}, false)

	// Assert
	assert.Equal(t, expected, gaps)
}

func TestWhenNoZonesAreGiven_ReturnsEmptyGapPerCapacity(t *testing.T) {
	t.Parallel()
	// Arrange
	expected := []neutralZone.Plans{nil, nil}

	// Act
	gaps := utils.AssignNeutralZonesToGaps(neutralZone.Plans{}, []int{1, 1}, false)

	// Assert
	assert.Equal(t, expected, gaps)
}
