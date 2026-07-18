package neutralZonePlan_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutralZone"
	"github.com/stretchr/testify/assert"
)

func TestWhenPlanListIsEmpty_ReturnsUnknownQuality(t *testing.T) {
	t.Parallel()
	// Arrange
	plans := neutralZone.Plans{}

	// Act
	quality := plans.GetQuality("A")

	// Assert
	assert.Equal(t, neutralZone.QualityUnknown, quality)
}

func TestWhenLabelIsNotFound_ReturnsUnknownQuality(t *testing.T) {
	t.Parallel()
	// Arrange
	plans := neutralZone.Plans{{Label: "A", Quality: neutralZone.QualityHigh}}

	// Act
	quality := plans.GetQuality("missing")

	// Assert
	assert.Equal(t, neutralZone.QualityUnknown, quality)
}

func TestWhenLabelIsFound_ReturnsThatPlanQuality(t *testing.T) {
	t.Parallel()
	// Arrange
	plans := neutralZone.Plans{
		{Label: "A", Quality: neutralZone.QualityHigh},
		{Label: "B", Quality: neutralZone.QualityLow},
	}

	// Act
	quality := plans.GetQuality("B")

	// Assert
	assert.Equal(t, neutralZone.QualityLow, quality)
}
