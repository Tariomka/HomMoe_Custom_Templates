package neutralZonePlan_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestWhenPlanListIsEmpty_ReturnsMediumQuality(t *testing.T) {
	// Arrange
	plans := models.NeutralZonePlans{}

	// Act
	quality := plans.GetQuality("A")

	// Assert
	assert.Equal(t, models.QualityMedium, quality)
}

func TestWhenLabelIsNotFound_ReturnsMediumQuality(t *testing.T) {
	// Arrange
	plans := models.NeutralZonePlans{{Label: "A", Quality: models.QualityHigh}}

	// Act
	quality := plans.GetQuality("missing")

	// Assert
	assert.Equal(t, models.QualityMedium, quality)
}

func TestWhenLabelIsFound_ReturnsThatPlanQuality(t *testing.T) {
	// Arrange
	plans := models.NeutralZonePlans{
		{Label: "A", Quality: models.QualityHigh},
		{Label: "B", Quality: models.QualityLow},
	}

	// Act
	quality := plans.GetQuality("B")

	// Assert
	assert.Equal(t, models.QualityLow, quality)
}
