package math_test

import (
	"math"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenValueHasMoreDecimalsThanPrecision_RoundsToPrecision(t *testing.T) {
	// Arrange
	value := 1.23456

	// Act
	actual := helpers.RoundWithPrecision(value, 2)

	// Assert
	assert.InDelta(t, 1.23, actual, test_helpers.Delta)
}

func TestWhenDecimalIsHalfway_RoundsAwayFromZero(t *testing.T) {
	// Arrange
	value := 0.125

	// Act
	actual := helpers.RoundWithPrecision(value, 2)

	// Assert
	assert.InDelta(t, 0.13, actual, test_helpers.Delta)
}

func TestWhenPrecisionIsZero_RoundsToNearestInteger(t *testing.T) {
	// Arrange
	value := gofakeit.Float64Range(-1000, 1000)

	// Act
	actual := helpers.RoundWithPrecision(value, 0)

	// Assert
	assert.InDelta(t, math.Round(value), actual, test_helpers.Delta)
}

func TestWhenValueIsFuzzed_MatchesIndependentRoundingFormula(t *testing.T) {
	// Arrange
	value := gofakeit.Float64Range(-1000, 1000)
	decimalPrecision := gofakeit.Number(1, 6)
	multiplier := math.Pow(10, float64(decimalPrecision))
	expected := math.Round(value*multiplier) / multiplier

	// Act
	actual := helpers.RoundWithPrecision(value, decimalPrecision)

	// Assert
	assert.InDelta(t, expected, actual, test_helpers.Delta)
}

func TestWhenValueIsNegative_RoundsToPrecision(t *testing.T) {
	// Arrange
	value := -2.71828

	// Act
	actual := helpers.RoundWithPrecision(value, 3)

	// Assert
	assert.InDelta(t, -2.718, actual, test_helpers.Delta)
}
