package misc_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones/utils"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenCapacitiesAreEmpty_ReturnsNoGaps(t *testing.T) {
	t.Parallel()
	// Arrange
	neutralZones := neutral_zone.Plans{
		{Label: "A", Quality: neutral_zone.QualityHigh, CastleCount: gofakeit.Number(0, 4)},
	}

	// Act
	gaps := utils.AssignNeutralZonesToGaps(neutralZones, []int{}, false)

	// Assert
	assert.Empty(t, gaps)
}

func TestWhenZonesFitExactly_AssignsStrongestZoneToLowestIndexedGap(t *testing.T) {
	t.Parallel()
	// Arrange
	strongZone := neutral_zone.Plan{Label: "A", Quality: neutral_zone.QualityHigh, CastleCount: gofakeit.Number(0, 4)}
	weakZone := neutral_zone.Plan{Label: "B", Quality: neutral_zone.QualityMedium, CastleCount: gofakeit.Number(0, 4)}
	neutralZones := neutral_zone.Plans{weakZone, strongZone}
	expected := []neutral_zone.Plans{{strongZone}, {weakZone}}

	// Act
	gaps := utils.AssignNeutralZonesToGaps(neutralZones, []int{1, 1}, false)

	// Assert
	assert.Equal(t, expected, gaps)
}

func TestWhenTotalCapacityIsExceeded_DropsWeakestZones(t *testing.T) {
	t.Parallel()
	// Arrange
	strongZone := neutral_zone.Plan{Label: "A", Quality: neutral_zone.QualityHigh, CastleCount: gofakeit.Number(0, 4)}
	weakZone := neutral_zone.Plan{Label: "B", Quality: neutral_zone.QualityLow, CastleCount: gofakeit.Number(0, 4)}
	neutralZones := neutral_zone.Plans{weakZone, strongZone}
	expected := []neutral_zone.Plans{{strongZone}}

	// Act
	gaps := utils.AssignNeutralZonesToGaps(neutralZones, []int{1}, false)

	// Assert
	assert.Equal(t, expected, gaps)
}

func TestWhenPreferInteriorIsTrue_AssignsFirstZoneToInteriorGap(t *testing.T) {
	t.Parallel()
	// Arrange
	zone := neutral_zone.Plan{Label: "A", Quality: neutral_zone.QualityHigh, CastleCount: gofakeit.Number(0, 4)}
	neutralZones := neutral_zone.Plans{zone}
	expected := []neutral_zone.Plans{nil, {zone}, nil}

	// Act
	gaps := utils.AssignNeutralZonesToGaps(neutralZones, []int{1, 1, 1}, true)

	// Assert
	assert.Equal(t, expected, gaps)
}

func TestWhenInteriorGapsAreFull_FallsBackToEdgeGaps(t *testing.T) {
	t.Parallel()
	// Arrange
	strongZone := neutral_zone.Plan{Label: "A", Quality: neutral_zone.QualityHigh, CastleCount: gofakeit.Number(0, 4)}
	mediumZone := neutral_zone.Plan{Label: "B", Quality: neutral_zone.QualityMedium, CastleCount: gofakeit.Number(0, 4)}
	weakZone := neutral_zone.Plan{Label: "C", Quality: neutral_zone.QualityLow, CastleCount: gofakeit.Number(0, 4)}
	neutralZones := neutral_zone.Plans{weakZone, mediumZone, strongZone}
	expected := []neutral_zone.Plans{{mediumZone}, {strongZone}, {weakZone}}

	// Act
	gaps := utils.AssignNeutralZonesToGaps(neutralZones, []int{1, 1, 1}, true)

	// Assert
	assert.Equal(t, expected, gaps)
}

func TestWhenOnlyTwoGapsExistWithPreferInterior_UsesEdgeGap(t *testing.T) {
	t.Parallel()
	// Arrange
	zone := neutral_zone.Plan{Label: "A", Quality: neutral_zone.QualityMedium, CastleCount: gofakeit.Number(0, 4)}
	neutralZones := neutral_zone.Plans{zone}
	expected := []neutral_zone.Plans{{zone}, nil}

	// Act
	gaps := utils.AssignNeutralZonesToGaps(neutralZones, []int{1, 1}, true)

	// Assert
	assert.Equal(t, expected, gaps)
}

func TestWhenGapLoadsAreTied_PrefersGapWithFewerZones(t *testing.T) {
	t.Parallel()
	// Arrange - one medium zone (score 2.0) balances two low zones (1.0 each),
	// so the fourth zone sees equal loads but fewer zones in the first gap.
	mediumZone := neutral_zone.Plan{Label: "A", Quality: neutral_zone.QualityMedium, CastleCount: 0}
	firstLowZone := neutral_zone.Plan{Label: "B", Quality: neutral_zone.QualityLow, CastleCount: 0}
	secondLowZone := neutral_zone.Plan{Label: "C", Quality: neutral_zone.QualityLow, CastleCount: 0}
	tieBreakerZone := neutral_zone.Plan{Label: "D", Quality: neutral_zone.QualityLow, CastleCount: 0}
	neutralZones := neutral_zone.Plans{mediumZone, firstLowZone, secondLowZone, tieBreakerZone}
	expected := []neutral_zone.Plans{{mediumZone, tieBreakerZone}, {firstLowZone, secondLowZone}}

	// Act
	gaps := utils.AssignNeutralZonesToGaps(neutralZones, []int{3, 3}, false)

	// Assert
	assert.Equal(t, expected, gaps)
}

func TestWhenNoZonesAreGiven_ReturnsEmptyGapPerCapacity(t *testing.T) {
	t.Parallel()
	// Arrange
	expected := []neutral_zone.Plans{nil, nil}

	// Act
	gaps := utils.AssignNeutralZonesToGaps(neutral_zone.Plans{}, []int{1, 1}, false)

	// Assert
	assert.Equal(t, expected, gaps)
}
