package neutralZonePlan_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutralZone"
	"github.com/stretchr/testify/assert"
)

func TestWhenPlansAreUnordered_SortsByQualityThenCastlesThenLabel(t *testing.T) {
	t.Parallel()
	// Arrange
	plans := neutralZone.Plans{
		{Label: "B", Quality: neutralZone.QualityLow, CastleCount: 0},
		{Label: "A", Quality: neutralZone.QualityLow, CastleCount: 0},
		{Label: "C", Quality: neutralZone.QualityHigh, CastleCount: 1},
		{Label: "D", Quality: neutralZone.QualityHigh, CastleCount: 2},
		{Label: "E", Quality: neutralZone.QualityMedium, CastleCount: 0},
	}
	expected := neutralZone.Plans{
		{Label: "D", Quality: neutralZone.QualityHigh, CastleCount: 2},
		{Label: "C", Quality: neutralZone.QualityHigh, CastleCount: 1},
		{Label: "E", Quality: neutralZone.QualityMedium, CastleCount: 0},
		{Label: "A", Quality: neutralZone.QualityLow, CastleCount: 0},
		{Label: "B", Quality: neutralZone.QualityLow, CastleCount: 0},
	}

	// Act
	sortedPlans := neutralZone.NewNeutralZonePlansSorted(plans)

	// Assert
	assert.Equal(t, expected, *sortedPlans)
}

func TestWhenSortedCopyIsCreated_LeavesInputSliceUntouched(t *testing.T) {
	t.Parallel()
	// Arrange
	plans := neutralZone.Plans{
		{Label: "B", Quality: neutralZone.QualityLow},
		{Label: "A", Quality: neutralZone.QualityHigh},
	}
	original := neutralZone.Plans{
		{Label: "B", Quality: neutralZone.QualityLow},
		{Label: "A", Quality: neutralZone.QualityHigh},
	}

	// Act
	neutralZone.NewNeutralZonePlansSorted(plans)

	// Assert
	assert.Equal(t, original, plans)
}
