package neutralZonePlan_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestWhenPlansArePresent_ReportsAny(t *testing.T) {
	t.Parallel()
	// Arrange
	plans := models.NeutralZonePlans{{Label: "A", Quality: models.QualityLow}}

	// Act
	anyPlans := plans.Any()

	// Assert
	assert.True(t, anyPlans)
}

func TestWhenPlanListIsEmpty_ReportsNone(t *testing.T) {
	t.Parallel()
	// Arrange
	plans := models.NeutralZonePlans{}

	// Act
	anyPlans := plans.Any()

	// Assert
	assert.False(t, anyPlans)
}
