package guardWeeklyIncrement_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_connections"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/stretchr/testify/assert"
)

func TestWhenCalled_ReturnsLabeledIncrementList(t *testing.T) {
	t.Parallel()
	// Arrange
	expected := []data.Tuple[string, float64]{
		data.NewTuple("Slow (5%)", 0.05),
		data.NewTuple("Normal (10%)", 0.10),
		data.NewTuple("Standard (15%)", 0.15),
		data.NewTuple("Fast (20%)", 0.20),
		data.NewTuple("Very Fast (25%)", 0.25),
	}

	// Act
	incrementList := common_connections.GetGuardWeeklyIncrementList()

	// Assert
	assert.Equal(t, expected, incrementList)
}
