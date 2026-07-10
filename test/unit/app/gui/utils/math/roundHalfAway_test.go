package math_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/utils"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
)

func TestWhenValueIsRounded_HalvesGoAwayFromZero(t *testing.T) {
	t.Parallel()
	// Arrange
	tests := []struct {
		name     string
		value    float64
		expected float64
	}{
		{name: "WhenFractionIsBelowHalf_RoundsDown", value: 2.4, expected: 2},
		{name: "WhenFractionIsExactlyHalf_RoundsUp", value: 2.5, expected: 3},
		{name: "WhenFractionIsAboveHalf_RoundsUp", value: 2.6, expected: 3},
		{name: "WhenNegativeFractionIsBelowHalf_RoundsTowardZero", value: -2.4, expected: -2},
		{name: "WhenNegativeFractionIsExactlyHalf_RoundsAwayFromZero", value: -2.5, expected: -3},
		{name: "WhenNegativeFractionIsAboveHalf_RoundsAwayFromZero", value: -2.6, expected: -3},
		{name: "WhenValueIsZero_ReturnsZero", value: 0, expected: 0},
		{name: "WhenValueIsWhole_ReturnsSameValue", value: 7, expected: 7},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			// Act
			result := utils.RoundHalfAway(testCase.value)

			// Assert
			assert.InDelta(t, testCase.expected, result, test_helpers.Delta)
		})
	}
}
