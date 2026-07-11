package neutralZonePlan_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestWhenTwoIndexesAreSwapped_ExchangesThosePlans(t *testing.T) {
	t.Parallel()
	// Arrange
	plans := models.NeutralZonePlans{
		{Label: "A", Quality: models.QualityLow},
		{Label: "B", Quality: models.QualityMedium},
		{Label: "C", Quality: models.QualityHigh},
	}
	expected := models.NeutralZonePlans{
		{Label: "C", Quality: models.QualityHigh},
		{Label: "B", Quality: models.QualityMedium},
		{Label: "A", Quality: models.QualityLow},
	}

	// Act
	plans.Swap(0, 2)

	// Assert
	assert.Equal(t, expected, plans)
}
