package neutralZonePlan_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/stretchr/testify/assert"
)

func TestWhenTwoIndexesAreSwapped_ExchangesThosePlans(t *testing.T) {
	t.Parallel()
	// Arrange
	plans := neutral_zone.Plans{
		{Label: "A", Quality: neutral_zone.QualityLow},
		{Label: "B", Quality: neutral_zone.QualityMedium},
		{Label: "C", Quality: neutral_zone.QualityHigh},
	}
	expected := neutral_zone.Plans{
		{Label: "C", Quality: neutral_zone.QualityHigh},
		{Label: "B", Quality: neutral_zone.QualityMedium},
		{Label: "A", Quality: neutral_zone.QualityLow},
	}

	// Act
	plans.Swap(0, 2)

	// Assert
	assert.Equal(t, expected, plans)
}
