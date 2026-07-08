package neutralZonePlan_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestWhenPlansAreUnordered_SortsByBalanceScoreDescendingThenLabel(t *testing.T) {
	// Arrange
	plans := models.NeutralZonePlans{
		{Label: "B", Quality: models.QualityLow, CastleCount: 0},    // score 1.0
		{Label: "A", Quality: models.QualityHigh, CastleCount: 2},   // score 3.3
		{Label: "C", Quality: models.QualityMedium, CastleCount: 1}, // score 2.15
	}
	expected := models.NeutralZonePlans{
		{Label: "A", Quality: models.QualityHigh, CastleCount: 2},
		{Label: "C", Quality: models.QualityMedium, CastleCount: 1},
		{Label: "B", Quality: models.QualityLow, CastleCount: 0},
	}

	// Act
	sortedPlans := models.NewNeutralZonePlansSortedByBalance(plans)

	// Assert
	assert.Equal(t, expected, *sortedPlans)
}

func TestWhenBalanceScoresAreEqual_BreaksTieByLabelAscending(t *testing.T) {
	// Arrange
	plans := models.NeutralZonePlans{
		{Label: "Z", Quality: models.QualityMedium, CastleCount: 1},
		{Label: "A", Quality: models.QualityMedium, CastleCount: 1},
	}
	expected := models.NeutralZonePlans{
		{Label: "A", Quality: models.QualityMedium, CastleCount: 1},
		{Label: "Z", Quality: models.QualityMedium, CastleCount: 1},
	}

	// Act
	sortedPlans := models.NewNeutralZonePlansSortedByBalance(plans)

	// Assert
	assert.Equal(t, expected, *sortedPlans)
}
