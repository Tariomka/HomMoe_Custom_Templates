package misc_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones/utils"
	"github.com/stretchr/testify/assert"
)

func TestGetEvenGapCapacities_TableDriven(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name      string
		gapCount  int
		itemCount int
		expected  []int
	}{
		{
			name:     "WhenGapCountIsZero_ReturnsNil",
			gapCount: 0, itemCount: 5,
			expected: nil,
		},
		{
			name:     "WhenGapCountIsNegative_ReturnsNil",
			gapCount: -2, itemCount: 5,
			expected: nil,
		},
		{
			name:     "WhenItemCountIsZero_ReturnsAllZeroCapacities",
			gapCount: 3, itemCount: 0,
			expected: []int{0, 0, 0},
		},
		{
			name:     "WhenItemCountIsNegative_ReturnsAllZeroCapacities",
			gapCount: 2, itemCount: -4,
			expected: []int{0, 0},
		},
		{
			name:     "WhenItemsDivideEvenly_SpreadsItemsEqually",
			gapCount: 2, itemCount: 4,
			expected: []int{2, 2},
		},
		{
			name:     "WhenItemsDivideUnevenly_GivesMiddleGapTheExtraItem",
			gapCount: 3, itemCount: 4,
			expected: []int{1, 2, 1},
		},
		{
			name:     "WhenItemsExceedGaps_SpreadsWithMiddleHeavy",
			gapCount: 3, itemCount: 7,
			expected: []int{2, 3, 2},
		},
		{
			name:     "WhenItemsFewerThanGaps_FillsOuterGapsFirst",
			gapCount: 3, itemCount: 2,
			expected: []int{1, 0, 1},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			// Arrange

			// Act
			capacities := utils.GetEvenGapCapacities(testCase.gapCount, testCase.itemCount)

			// Assert
			assert.Equal(t, testCase.expected, capacities)
		})
	}
}
