package math_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/utils"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenValueIsInsideUnitRange_LinearInterpolationIsReturned(t *testing.T) {
	// Arrange
	tests := []struct {
		name     string
		value    float32
		low      float32
		high     float32
		expected float32
	}{
		{name: "WhenValueIsZero_ReturnsLow", value: 0, low: 2, high: 10, expected: 2},
		{name: "WhenValueIsOne_ReturnsHigh", value: 1, low: 2, high: 10, expected: 10},
		{name: "WhenValueIsMidway_ReturnsMidpoint", value: 0.5, low: 2, high: 10, expected: 6},
		{name: "WhenRangeIsNegative_InterpolatesDownward", value: 0.25, low: 0, high: -8, expected: -2},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			// Act
			result := utils.Denormalize(testCase.value, testCase.low, testCase.high)

			// Assert
			assert.InDelta(t, testCase.expected, result, test_helpers.Delta)
		})
	}
}

func TestWhenValueIsBelowZero_ClampsToLow(t *testing.T) {
	// Arrange
	low := gofakeit.Float32Range(-100, 0)
	high := gofakeit.Float32Range(1, 100)

	// Act
	result := utils.Denormalize(-0.5, low, high)

	// Assert
	assert.InDelta(t, low, result, test_helpers.Delta)
}

func TestWhenValueIsAboveOne_ClampsToHigh(t *testing.T) {
	// Arrange
	low := gofakeit.Float32Range(-100, 0)
	high := gofakeit.Float32Range(1, 100)

	// Act
	result := utils.Denormalize(1.5, low, high)

	// Assert - low + 1*(high-low) may differ from high by a float32 rounding step
	assert.InDelta(t, high, result, test_helpers.Delta)
}
