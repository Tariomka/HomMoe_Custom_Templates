package math_test

import (
	stdmath "math"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenValueHasMoreDecimalsThanPrecision_RoundsToPrecision(t *testing.T) {
	// Arrange
	value := 1.23456

	// Act
	rounded := helpers.RoundWithPrecision(value, 2)

	// Assert
	assert.Equal(t, 1.23, rounded)
}

func TestWhenDecimalIsHalfway_RoundsAwayFromZero(t *testing.T) {
	// Arrange
	value := 0.125

	// Act
	rounded := helpers.RoundWithPrecision(value, 2)

	// Assert
	assert.Equal(t, 0.13, rounded)
}

func TestWhenPrecisionIsZero_RoundsToNearestInteger(t *testing.T) {
	// Arrange
	value := gofakeit.Float64Range(-1000, 1000)

	// Act
	rounded := helpers.RoundWithPrecision(value, 0)

	// Assert
	assert.Equal(t, stdmath.Round(value), rounded)
}

func TestWhenValueIsFuzzed_MatchesIndependentRoundingFormula(t *testing.T) {
	// Arrange
	value := gofakeit.Float64Range(-1000, 1000)
	decimalPrecision := gofakeit.Number(1, 6)
	multiplier := stdmath.Pow(10, float64(decimalPrecision))
	expected := stdmath.Round(value*multiplier) / multiplier

	// Act
	rounded := helpers.RoundWithPrecision(value, decimalPrecision)

	// Assert
	assert.Equal(t, expected, rounded)
}

func TestWhenValueIsNegative_RoundsToPrecision(t *testing.T) {
	// Arrange
	value := -2.71828

	// Act
	rounded := helpers.RoundWithPrecision(value, 3)

	// Assert
	assert.Equal(t, -2.718, rounded)
}
