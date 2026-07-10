package math_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/utils"
	"github.com/stretchr/testify/assert"
)

func TestWhenSliderValueIsSnapped_NearestIntegerInRangeIsReturned(t *testing.T) {
	t.Parallel()
	// Arrange
	tests := []struct {
		name     string
		value    float32
		low      int
		high     int
		expected int
	}{
		{name: "WhenSliderIsAtZero_ReturnsLowBound", value: 0, low: 3, high: 7, expected: 3},
		{name: "WhenSliderIsAtOne_ReturnsHighBound", value: 1, low: 3, high: 7, expected: 7},
		{name: "WhenSliderIsMidway_ReturnsMiddleInteger", value: 0.5, low: 0, high: 10, expected: 5},
		{name: "WhenSliderIsMidwayOverNegativeRange_ReturnsZero", value: 0.5, low: -5, high: 5, expected: 0},
		{name: "WhenSliderIsBelowZero_ClampsToLowBound", value: -0.4, low: 2, high: 9, expected: 2},
		{name: "WhenSliderIsAboveOne_ClampsToHighBound", value: 1.4, low: 2, high: 9, expected: 9},
		{name: "WhenSnapFallsBetweenIntegers_RoundsToNearest", value: 0.26, low: 0, high: 10, expected: 3},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			// Act
			result := utils.RoundedRange(testCase.value, testCase.low, testCase.high)

			// Assert
			assert.Equal(t, testCase.expected, result)
		})
	}
}
