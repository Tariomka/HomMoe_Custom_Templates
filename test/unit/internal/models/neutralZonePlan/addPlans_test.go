package neutralZonePlan_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestWhenMultiplePlansAreAdded_AppendsAllInOrder(t *testing.T) {
	// Arrange
	plans := models.NeutralZonePlans{{Label: "A", Quality: models.QualityLow}}
	added := []models.NeutralZonePlan{
		{Label: "B", Quality: models.QualityMedium, CastleCount: 1},
		{Label: "C", Quality: models.QualityHigh, CastleCount: 2},
	}
	expected := models.NeutralZonePlans{
		{Label: "A", Quality: models.QualityLow},
		{Label: "B", Quality: models.QualityMedium, CastleCount: 1},
		{Label: "C", Quality: models.QualityHigh, CastleCount: 2},
	}

	// Act
	plans.AddPlans(added...)

	// Assert
	assert.Equal(t, expected, plans)
}
