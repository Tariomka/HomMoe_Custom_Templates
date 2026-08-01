package guardWeeklyIncrement_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_connections"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestWhenCalled_ReturnsWeeklyIncrementValues(t *testing.T) {
	t.Parallel()
	// Arrange
	expected := models.GuardWeeklyIncrement{
		Slow:     0.05,
		Normal:   0.10,
		Standard: 0.15,
		Fast:     0.20,
		VeryFast: 0.25,
	}

	// Act
	increments := common_connections.GetGuardWeeklyIncrements()

	// Assert
	assert.Equal(t, expected, increments)
}
