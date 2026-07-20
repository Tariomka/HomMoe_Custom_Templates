package neutralZonePlan_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/stretchr/testify/assert"
)

func TestWhenPlansArePresent_ReportsAny(t *testing.T) {
	t.Parallel()
	// Arrange
	plans := neutral_zone.Plans{{Label: "A", Quality: neutral_zone.QualityLow}}

	// Act
	anyPlans := plans.Any()

	// Assert
	assert.True(t, anyPlans)
}

func TestWhenPlanListIsEmpty_ReportsNone(t *testing.T) {
	t.Parallel()
	// Arrange
	plans := neutral_zone.Plans{}

	// Act
	anyPlans := plans.Any()

	// Assert
	assert.False(t, anyPlans)
}
