package math_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/utils"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenValueIsInsideRange_RatioIsReturned(t *testing.T) {
	// Arrange
	tests := []struct {
		name     string
		value    float32
		low      float32
		high     float32
		expected float32
	}{
		{name: "WhenValueIsAtLow_ReturnsZero", value: 2, low: 2, high: 10, expected: 0},
		{name: "WhenValueIsAtHigh_ReturnsOne", value: 10, low: 2, high: 10, expected: 1},
		{name: "WhenValueIsMidway_ReturnsHalf", value: 6, low: 2, high: 10, expected: 0.5},
		{name: "WhenValueIsAtQuarter_ReturnsQuarterRatio", value: 4, low: 2, high: 10, expected: 0.25},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			// Act
			result := utils.Normalize(testCase.value, testCase.low, testCase.high)

			// Assert
			assert.InDelta(t, testCase.expected, result, 0.0001)
		})
	}
}

func TestWhenLowEqualsHigh_ZeroIsReturned(t *testing.T) {
	// Arrange
	bound := gofakeit.Float32Range(-100, 100)

	// Act
	result := utils.Normalize(gofakeit.Float32Range(-100, 100), bound, bound)

	// Assert
	assert.Equal(t, float32(0), result)
}

func TestWhenValueIsBelowRange_ClampsToZero(t *testing.T) {
	// Arrange
	low := gofakeit.Float32Range(0, 10)
	high := gofakeit.Float32Range(20, 100)

	// Act
	result := utils.Normalize(low-5, low, high)

	// Assert
	assert.Equal(t, float32(0), result)
}

func TestWhenValueIsAboveRange_ClampsToOne(t *testing.T) {
	// Arrange
	low := gofakeit.Float32Range(0, 10)
	high := gofakeit.Float32Range(20, 100)

	// Act
	result := utils.Normalize(high+5, low, high)

	// Assert
	assert.Equal(t, float32(1), result)
}
