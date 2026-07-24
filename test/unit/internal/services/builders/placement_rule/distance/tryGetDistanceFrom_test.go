package distance_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/placement_rule"
	"github.com/stretchr/testify/assert"
)

type distanceLookupResult struct {
	Distance placement_rule.Distance
	Found    bool
}

func TestWhenLabelIsResolved_ReturnsMatchingPresetAndFoundFlag(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name     string
		label    string
		expected distanceLookupResult
	}{
		{
			name:     "WhenLabelIsNextTo_ReturnsNextToPreset",
			label:    "Next To",
			expected: distanceLookupResult{Distance: placement_rule.DistanceNextTo, Found: true},
		},
		{
			name:  "WhenLabelIsNear_ReturnsNearPreset",
			label: "Near",
			expected: distanceLookupResult{
				Distance: placement_rule.Distance{Min: 0.075, Max: 0.35},
				Found:    true,
			},
		},
		{
			name:     "WhenLabelIsMedium_ReturnsMediumPreset",
			label:    "Medium",
			expected: distanceLookupResult{Distance: placement_rule.DistanceMedium, Found: true},
		},
		{
			name:     "WhenLabelIsFar_ReturnsFarPreset",
			label:    "Far",
			expected: distanceLookupResult{Distance: placement_rule.DistanceFar, Found: true},
		},
		{
			name:     "WhenLabelIsVeryFar_ReturnsVeryFarPreset",
			label:    "Very Far",
			expected: distanceLookupResult{Distance: placement_rule.DistanceVeryFar, Found: true},
		},
		{
			name:     "WhenLabelHasSurroundingWhitespace_StillResolvesPreset",
			label:    "  Near  ",
			expected: distanceLookupResult{Distance: placement_rule.DistanceNear, Found: true},
		},
		{
			name:     "WhenLabelIsEmpty_ReturnsNotFound",
			label:    "",
			expected: distanceLookupResult{Found: false},
		},
		{
			name:     "WhenLabelIsUnknown_ReturnsNotFound",
			label:    "Somewhere Else",
			expected: distanceLookupResult{Found: false},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			// Arrange & Act
			actualDistance, actualFound := placement_rule.TryGetDistanceFrom(testCase.label)

			// Assert
			assert.Equal(t, testCase.expected, distanceLookupResult{Distance: actualDistance, Found: actualFound})
		})
	}
}
