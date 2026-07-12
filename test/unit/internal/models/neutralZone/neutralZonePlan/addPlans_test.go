package neutralZonePlan_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutralZone"
	"github.com/stretchr/testify/assert"
)

func TestWhenMultiplePlansAreAdded_AppendsAllInOrder(t *testing.T) {
	t.Parallel()
	// Arrange
	plans := neutralZone.Plans{{Label: "A", Quality: neutralZone.QualityLow}}
	added := []neutralZone.Plan{
		{Label: "B", Quality: neutralZone.QualityMedium, CastleCount: 1},
		{Label: "C", Quality: neutralZone.QualityHigh, CastleCount: 2},
	}
	expected := neutralZone.Plans{
		{Label: "A", Quality: neutralZone.QualityLow},
		{Label: "B", Quality: neutralZone.QualityMedium, CastleCount: 1},
		{Label: "C", Quality: neutralZone.QualityHigh, CastleCount: 2},
	}

	// Act
	plans.AddPlans(added...)

	// Assert
	assert.Equal(t, expected, plans)
}
