package neutralZonePlan_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/stretchr/testify/assert"
)

func TestWhenPlansAreUnordered_SortsByBalanceScoreDescendingThenLabel(t *testing.T) {
	t.Parallel()
	// Arrange
	plans := neutral_zone.Plans{
		{Label: "B", Quality: neutral_zone.QualityLow, CastleCount: 0},    // score 1.0
		{Label: "A", Quality: neutral_zone.QualityHigh, CastleCount: 2},   // score 3.3
		{Label: "C", Quality: neutral_zone.QualityMedium, CastleCount: 1}, // score 2.15
	}
	expected := neutral_zone.Plans{
		{Label: "A", Quality: neutral_zone.QualityHigh, CastleCount: 2},
		{Label: "C", Quality: neutral_zone.QualityMedium, CastleCount: 1},
		{Label: "B", Quality: neutral_zone.QualityLow, CastleCount: 0},
	}

	// Act
	sortedPlans := neutral_zone.NewNeutralZonePlansSortedByBalance(plans)

	// Assert
	assert.Equal(t, expected, *sortedPlans)
}

func TestWhenBalanceScoresAreEqual_BreaksTieByLabelAscending(t *testing.T) {
	t.Parallel()
	// Arrange
	plans := neutral_zone.Plans{
		{Label: "Z", Quality: neutral_zone.QualityMedium, CastleCount: 1},
		{Label: "A", Quality: neutral_zone.QualityMedium, CastleCount: 1},
	}
	expected := neutral_zone.Plans{
		{Label: "A", Quality: neutral_zone.QualityMedium, CastleCount: 1},
		{Label: "Z", Quality: neutral_zone.QualityMedium, CastleCount: 1},
	}

	// Act
	sortedPlans := neutral_zone.NewNeutralZonePlansSortedByBalance(plans)

	// Assert
	assert.Equal(t, expected, *sortedPlans)
}
