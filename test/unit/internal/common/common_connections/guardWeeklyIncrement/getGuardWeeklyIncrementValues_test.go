package guardWeeklyIncrement_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_connections"
	"github.com/stretchr/testify/assert"
)

func TestWhenCalled_ReturnsIncrementValuesInOrder(t *testing.T) {
	t.Parallel()
	// Arrange
	expected := []float64{0.05, 0.10, 0.15, 0.20, 0.25}

	// Act
	values := common_connections.GetGuardWeeklyIncrementValues()

	// Assert
	assert.Equal(t, expected, values)
}
