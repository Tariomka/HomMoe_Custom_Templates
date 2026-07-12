package neutralZonePlan_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutralZone"
	"github.com/stretchr/testify/assert"
)

func TestWhenPlansAreUnordered_SortsByBalanceScoreAscendingThenLabel(t *testing.T) {
	t.Parallel()
	// Arrange
	plans := neutralZone.Plans{
		{Label: "A", Quality: neutralZone.QualityHigh, CastleCount: 2},   // score 3.3
		{Label: "B", Quality: neutralZone.QualityLow, CastleCount: 0},    // score 1.0
		{Label: "C", Quality: neutralZone.QualityMedium, CastleCount: 1}, // score 2.15
	}
	expected := neutralZone.Plans{
		{Label: "B", Quality: neutralZone.QualityLow, CastleCount: 0},
		{Label: "C", Quality: neutralZone.QualityMedium, CastleCount: 1},
		{Label: "A", Quality: neutralZone.QualityHigh, CastleCount: 2},
	}

	// Act
	plans.SortByBalanceScoreAscending()

	// Assert
	assert.Equal(t, expected, plans)
}

func TestWhenBalanceScoresAreEqualDuringAscendingSort_BreaksTieByLabelAscending(t *testing.T) {
	t.Parallel()
	// Arrange
	plans := neutralZone.Plans{
		{Label: "Z", Quality: neutralZone.QualityLow, CastleCount: 1},
		{Label: "A", Quality: neutralZone.QualityLow, CastleCount: 1},
	}
	expected := neutralZone.Plans{
		{Label: "A", Quality: neutralZone.QualityLow, CastleCount: 1},
		{Label: "Z", Quality: neutralZone.QualityLow, CastleCount: 1},
	}

	// Act
	plans.SortByBalanceScoreAscending()

	// Assert
	assert.Equal(t, expected, plans)
}
