package neutralZonePlan_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/stretchr/testify/assert"
)

func TestWhenMultiplePlansAreAdded_AppendsAllInOrder(t *testing.T) {
	t.Parallel()
	// Arrange
	plans := neutral_zone.Plans{{Label: "A", Quality: neutral_zone.QualityLow}}
	added := []neutral_zone.Plan{
		{Label: "B", Quality: neutral_zone.QualityMedium, CastleCount: 1},
		{Label: "C", Quality: neutral_zone.QualityHigh, CastleCount: 2},
	}
	expected := neutral_zone.Plans{
		{Label: "A", Quality: neutral_zone.QualityLow},
		{Label: "B", Quality: neutral_zone.QualityMedium, CastleCount: 1},
		{Label: "C", Quality: neutral_zone.QualityHigh, CastleCount: 2},
	}

	// Act
	plans.AddPlans(added...)

	// Assert
	assert.Equal(t, expected, plans)
}
