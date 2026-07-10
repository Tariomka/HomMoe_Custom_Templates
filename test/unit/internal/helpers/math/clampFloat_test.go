package math_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenFloatValueIsBelowLowest_ReturnsLowest(t *testing.T) {
	t.Parallel()
	// Arrange
	lowest := gofakeit.Float64Range(0, 100)
	highest := lowest + gofakeit.Float64Range(1, 100)
	value := lowest - gofakeit.Float64Range(1, 100)

	// Act
	actual := helpers.ClampFloat(value, lowest, highest)

	// Assert
	assert.InDelta(t, lowest, actual, test_helpers.Delta)
}

func TestWhenFloatValueIsAboveHighest_ReturnsHighest(t *testing.T) {
	t.Parallel()
	// Arrange
	lowest := gofakeit.Float64Range(0, 100)
	highest := lowest + gofakeit.Float64Range(1, 100)
	value := highest + gofakeit.Float64Range(1, 100)

	// Act
	actual := helpers.ClampFloat(value, lowest, highest)

	// Assert
	assert.InDelta(t, highest, actual, test_helpers.Delta)
}

func TestWhenFloatValueIsWithinBounds_ReturnsValueUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	lowest := gofakeit.Float64Range(0, 100)
	highest := lowest + gofakeit.Float64Range(2, 100)
	value := gofakeit.Float64Range(lowest, highest)

	// Act
	actual := helpers.ClampFloat(value, lowest, highest)

	// Assert
	assert.InDelta(t, value, actual, test_helpers.Delta)
}

func TestWhenFloatValueEqualsLowest_ReturnsValueUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	lowest := gofakeit.Float64Range(0, 100)
	highest := lowest + gofakeit.Float64Range(1, 100)

	// Act
	actual := helpers.ClampFloat(lowest, lowest, highest)

	// Assert
	assert.InDelta(t, lowest, actual, test_helpers.Delta)
}

func TestWhenFloatValueEqualsHighest_ReturnsValueUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	lowest := gofakeit.Float64Range(0, 100)
	highest := lowest + gofakeit.Float64Range(1, 100)

	// Act
	actual := helpers.ClampFloat(highest, lowest, highest)

	// Assert
	assert.InDelta(t, highest, actual, test_helpers.Delta)
}
