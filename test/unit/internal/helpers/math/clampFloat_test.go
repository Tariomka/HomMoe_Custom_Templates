package math_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenFloatValueIsBelowLowest_ReturnsLowest(t *testing.T) {
	// Arrange
	lowest := gofakeit.Float64Range(0, 100)
	highest := lowest + gofakeit.Float64Range(1, 100)
	value := lowest - gofakeit.Float64Range(1, 100)

	// Act
	clamped := helpers.ClampFloat(value, lowest, highest)

	// Assert
	assert.Equal(t, lowest, clamped)
}

func TestWhenFloatValueIsAboveHighest_ReturnsHighest(t *testing.T) {
	// Arrange
	lowest := gofakeit.Float64Range(0, 100)
	highest := lowest + gofakeit.Float64Range(1, 100)
	value := highest + gofakeit.Float64Range(1, 100)

	// Act
	clamped := helpers.ClampFloat(value, lowest, highest)

	// Assert
	assert.Equal(t, highest, clamped)
}

func TestWhenFloatValueIsWithinBounds_ReturnsValueUnchanged(t *testing.T) {
	// Arrange
	lowest := gofakeit.Float64Range(0, 100)
	highest := lowest + gofakeit.Float64Range(2, 100)
	value := gofakeit.Float64Range(lowest, highest)

	// Act
	clamped := helpers.ClampFloat(value, lowest, highest)

	// Assert
	assert.Equal(t, value, clamped)
}

func TestWhenFloatValueEqualsLowest_ReturnsValueUnchanged(t *testing.T) {
	// Arrange
	lowest := gofakeit.Float64Range(0, 100)
	highest := lowest + gofakeit.Float64Range(1, 100)

	// Act
	clamped := helpers.ClampFloat(lowest, lowest, highest)

	// Assert
	assert.Equal(t, lowest, clamped)
}

func TestWhenFloatValueEqualsHighest_ReturnsValueUnchanged(t *testing.T) {
	// Arrange
	lowest := gofakeit.Float64Range(0, 100)
	highest := lowest + gofakeit.Float64Range(1, 100)

	// Act
	clamped := helpers.ClampFloat(highest, lowest, highest)

	// Assert
	assert.Equal(t, highest, clamped)
}
