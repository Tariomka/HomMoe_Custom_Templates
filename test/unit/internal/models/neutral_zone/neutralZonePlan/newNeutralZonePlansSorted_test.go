package neutralZonePlan_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/stretchr/testify/assert"
)

func TestWhenPlansAreUnordered_SortsByQualityThenCastlesThenLabel(t *testing.T) {
	t.Parallel()
	// Arrange
	plans := neutral_zone.Plans{
		{Label: "B", Quality: neutral_zone.QualityLow, CastleCount: 0},
		{Label: "A", Quality: neutral_zone.QualityLow, CastleCount: 0},
		{Label: "C", Quality: neutral_zone.QualityHigh, CastleCount: 1},
		{Label: "D", Quality: neutral_zone.QualityHigh, CastleCount: 2},
		{Label: "E", Quality: neutral_zone.QualityMedium, CastleCount: 0},
	}
	expected := neutral_zone.Plans{
		{Label: "D", Quality: neutral_zone.QualityHigh, CastleCount: 2},
		{Label: "C", Quality: neutral_zone.QualityHigh, CastleCount: 1},
		{Label: "E", Quality: neutral_zone.QualityMedium, CastleCount: 0},
		{Label: "A", Quality: neutral_zone.QualityLow, CastleCount: 0},
		{Label: "B", Quality: neutral_zone.QualityLow, CastleCount: 0},
	}

	// Act
	sortedPlans := neutral_zone.NewNeutralZonePlansSorted(plans)

	// Assert
	assert.Equal(t, expected, *sortedPlans)
}

func TestWhenSortedCopyIsCreated_LeavesInputSliceUntouched(t *testing.T) {
	t.Parallel()
	// Arrange
	plans := neutral_zone.Plans{
		{Label: "B", Quality: neutral_zone.QualityLow},
		{Label: "A", Quality: neutral_zone.QualityHigh},
	}
	original := neutral_zone.Plans{
		{Label: "B", Quality: neutral_zone.QualityLow},
		{Label: "A", Quality: neutral_zone.QualityHigh},
	}

	// Act
	neutral_zone.NewNeutralZonePlansSorted(plans)

	// Assert
	assert.Equal(t, original, plans)
}
