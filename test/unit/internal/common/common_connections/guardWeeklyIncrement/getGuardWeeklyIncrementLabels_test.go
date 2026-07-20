package guardWeeklyIncrement_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_connections"
	"github.com/stretchr/testify/assert"
)

func TestWhenCalled_ReturnsIncrementLabelsInOrder(t *testing.T) {
	t.Parallel()
	// Arrange
	expected := []string{"Slow (5%)", "Normal (10%)", "Standard (15%)", "Fast (20%)", "Very Fast (25%)"}

	// Act
	labels := common_connections.GetGuardWeeklyIncrementLabels()

	// Assert
	assert.Equal(t, expected, labels)
}
