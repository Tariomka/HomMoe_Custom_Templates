package neutralZonePlan_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutralZone"
	"github.com/stretchr/testify/assert"
)

func TestWhenTwoIndexesAreSwapped_ExchangesThosePlans(t *testing.T) {
	t.Parallel()
	// Arrange
	plans := neutralZone.Plans{
		{Label: "A", Quality: neutralZone.QualityLow},
		{Label: "B", Quality: neutralZone.QualityMedium},
		{Label: "C", Quality: neutralZone.QualityHigh},
	}
	expected := neutralZone.Plans{
		{Label: "C", Quality: neutralZone.QualityHigh},
		{Label: "B", Quality: neutralZone.QualityMedium},
		{Label: "A", Quality: neutralZone.QualityLow},
	}

	// Act
	plans.Swap(0, 2)

	// Assert
	assert.Equal(t, expected, plans)
}
