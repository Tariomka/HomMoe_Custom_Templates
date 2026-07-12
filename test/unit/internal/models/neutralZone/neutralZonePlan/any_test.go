package neutralZonePlan_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutralZone"
	"github.com/stretchr/testify/assert"
)

func TestWhenPlansArePresent_ReportsAny(t *testing.T) {
	t.Parallel()
	// Arrange
	plans := neutralZone.Plans{{Label: "A", Quality: neutralZone.QualityLow}}

	// Act
	anyPlans := plans.Any()

	// Assert
	assert.True(t, anyPlans)
}

func TestWhenPlanListIsEmpty_ReportsNone(t *testing.T) {
	t.Parallel()
	// Arrange
	plans := neutralZone.Plans{}

	// Act
	anyPlans := plans.Any()

	// Assert
	assert.False(t, anyPlans)
}
