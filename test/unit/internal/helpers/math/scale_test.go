package math_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenProductIsPositive_ReturnsTruncatedProduct(t *testing.T) {
	t.Parallel()
	// Arrange
	value := gofakeit.Float64Range(1, 500)
	multiplier := gofakeit.Float64Range(0.1, 10)
	expected := int(value * multiplier)

	// Act
	scaled := helpers.Scale(value, multiplier)

	// Assert
	assert.Equal(t, expected, scaled)
}

func TestWhenProductIsNegative_ReturnsZero(t *testing.T) {
	t.Parallel()
	// Arrange
	value := gofakeit.Float64Range(1, 500)
	multiplier := -gofakeit.Float64Range(0.1, 10)

	// Act
	scaled := helpers.Scale(value, multiplier)

	// Assert
	assert.Equal(t, 0, scaled)
}

func TestWhenProductIsFractionBelowOne_ReturnsZero(t *testing.T) {
	t.Parallel()
	// Arrange
	value := 0.4
	multiplier := 2.0

	// Act
	scaled := helpers.Scale(value, multiplier)

	// Assert
	assert.Equal(t, 0, scaled)
}

func TestWhenValueIsZero_ReturnsZero(t *testing.T) {
	t.Parallel()
	// Arrange
	multiplier := gofakeit.Float64Range(0.1, 10)

	// Act
	scaled := helpers.Scale(0, multiplier)

	// Assert
	assert.Equal(t, 0, scaled)
}
