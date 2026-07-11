package neutralZonePlan_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestWhenPlansAreUnordered_SortsByQualityThenCastlesThenLabel(t *testing.T) {
	t.Parallel()
	// Arrange
	plans := models.NeutralZonePlans{
		{Label: "B", Quality: models.QualityLow, CastleCount: 0},
		{Label: "A", Quality: models.QualityLow, CastleCount: 0},
		{Label: "C", Quality: models.QualityHigh, CastleCount: 1},
		{Label: "D", Quality: models.QualityHigh, CastleCount: 2},
		{Label: "E", Quality: models.QualityMedium, CastleCount: 0},
	}
	expected := models.NeutralZonePlans{
		{Label: "D", Quality: models.QualityHigh, CastleCount: 2},
		{Label: "C", Quality: models.QualityHigh, CastleCount: 1},
		{Label: "E", Quality: models.QualityMedium, CastleCount: 0},
		{Label: "A", Quality: models.QualityLow, CastleCount: 0},
		{Label: "B", Quality: models.QualityLow, CastleCount: 0},
	}

	// Act
	sortedPlans := models.NewNeutralZonePlansSorted(plans)

	// Assert
	assert.Equal(t, expected, *sortedPlans)
}

func TestWhenSortedCopyIsCreated_LeavesInputSliceUntouched(t *testing.T) {
	t.Parallel()
	// Arrange
	plans := models.NeutralZonePlans{
		{Label: "B", Quality: models.QualityLow},
		{Label: "A", Quality: models.QualityHigh},
	}
	original := models.NeutralZonePlans{
		{Label: "B", Quality: models.QualityLow},
		{Label: "A", Quality: models.QualityHigh},
	}

	// Act
	models.NewNeutralZonePlansSorted(plans)

	// Assert
	assert.Equal(t, original, plans)
}
