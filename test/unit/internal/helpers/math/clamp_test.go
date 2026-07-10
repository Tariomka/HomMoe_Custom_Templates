package math_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenValueIsBelowLowest_ReturnsLowest(t *testing.T) {
	t.Parallel()
	// Arrange
	lowest := gofakeit.Number(0, 100)
	highest := lowest + gofakeit.Number(1, 100)
	value := lowest - gofakeit.Number(1, 100)

	// Act
	clamped := helpers.Clamp(value, lowest, highest)

	// Assert
	assert.Equal(t, lowest, clamped)
}

func TestWhenValueIsAboveHighest_ReturnsHighest(t *testing.T) {
	t.Parallel()
	// Arrange
	lowest := gofakeit.Number(0, 100)
	highest := lowest + gofakeit.Number(1, 100)
	value := highest + gofakeit.Number(1, 100)

	// Act
	clamped := helpers.Clamp(value, lowest, highest)

	// Assert
	assert.Equal(t, highest, clamped)
}

func TestWhenValueIsWithinBounds_ReturnsValueUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	lowest := gofakeit.Number(0, 100)
	highest := lowest + gofakeit.Number(2, 100)
	value := gofakeit.Number(lowest, highest)

	// Act
	clamped := helpers.Clamp(value, lowest, highest)

	// Assert
	assert.Equal(t, value, clamped)
}

func TestWhenValueEqualsLowest_ReturnsValueUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	lowest := gofakeit.Number(0, 100)
	highest := lowest + gofakeit.Number(1, 100)

	// Act
	clamped := helpers.Clamp(lowest, lowest, highest)

	// Assert
	assert.Equal(t, lowest, clamped)
}

func TestWhenValueEqualsHighest_ReturnsValueUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	lowest := gofakeit.Number(0, 100)
	highest := lowest + gofakeit.Number(1, 100)

	// Act
	clamped := helpers.Clamp(highest, lowest, highest)

	// Assert
	assert.Equal(t, highest, clamped)
}
