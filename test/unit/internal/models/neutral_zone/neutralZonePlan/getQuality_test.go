package neutralZonePlan_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/stretchr/testify/assert"
)

func TestWhenPlanListIsEmpty_ReturnsUnknownQuality(t *testing.T) {
	t.Parallel()
	// Arrange
	plans := neutral_zone.Plans{}

	// Act
	quality := plans.GetQuality("A")

	// Assert
	assert.Equal(t, neutral_zone.QualityUnknown, quality)
}

func TestWhenLabelIsNotFound_ReturnsUnknownQuality(t *testing.T) {
	t.Parallel()
	// Arrange
	plans := neutral_zone.Plans{{Label: "A", Quality: neutral_zone.QualityHigh}}

	// Act
	quality := plans.GetQuality("missing")

	// Assert
	assert.Equal(t, neutral_zone.QualityUnknown, quality)
}

func TestWhenLabelIsFound_ReturnsThatPlanQuality(t *testing.T) {
	t.Parallel()
	// Arrange
	plans := neutral_zone.Plans{
		{Label: "A", Quality: neutral_zone.QualityHigh},
		{Label: "B", Quality: neutral_zone.QualityLow},
	}

	// Act
	quality := plans.GetQuality("B")

	// Assert
	assert.Equal(t, neutral_zone.QualityLow, quality)
}
