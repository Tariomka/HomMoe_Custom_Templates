package misc_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones/utils"
	"github.com/stretchr/testify/assert"
)

func TestGetEvenGapCapacities_TableDriven(t *testing.T) {
	testCases := []struct {
		name          string
		gapCount      int
		itemCount     int
		minimumPerGap int
		expected      []int
	}{
		{
			name:     "WhenGapCountIsZero_ReturnsNil",
			gapCount: 0, itemCount: 5, minimumPerGap: 0,
			expected: nil,
		},
		{
			name:     "WhenGapCountIsNegative_ReturnsNil",
			gapCount: -2, itemCount: 5, minimumPerGap: 1,
			expected: nil,
		},
		{
			name:     "WhenItemCountIsZero_ReturnsAllZeroCapacities",
			gapCount: 3, itemCount: 0, minimumPerGap: 1,
			expected: []int{0, 0, 0},
		},
		{
			name:     "WhenItemCountIsNegative_ReturnsAllZeroCapacities",
			gapCount: 2, itemCount: -4, minimumPerGap: 0,
			expected: []int{0, 0},
		},
		{
			name:     "WhenItemsDivideEvenly_SpreadsItemsEqually",
			gapCount: 2, itemCount: 4, minimumPerGap: 0,
			expected: []int{2, 2},
		},
		{
			name:     "WhenItemsDivideUnevenly_GivesMiddleGapTheExtraItem",
			gapCount: 3, itemCount: 4, minimumPerGap: 0,
			expected: []int{1, 2, 1},
		},
		{
			name:     "WhenMinimumFitsEveryGap_ReservesMinimumThenSpreadsRemainder",
			gapCount: 3, itemCount: 7, minimumPerGap: 2,
			expected: []int{2, 3, 2},
		},
		{
			name:     "WhenItemsCannotCoverMinimumForEveryGap_IgnoresMinimum",
			gapCount: 3, itemCount: 2, minimumPerGap: 1,
			expected: []int{1, 0, 1},
		},
		{
			name:     "WhenMinimumIsNegative_TreatsItAsZero",
			gapCount: 2, itemCount: 2, minimumPerGap: -3,
			expected: []int{1, 1},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// Arrange - inputs come from the table entry

			// Act
			capacities := utils.GetEvenGapCapacities(testCase.gapCount, testCase.itemCount, testCase.minimumPerGap)

			// Assert
			assert.Equal(t, testCase.expected, capacities)
		})
	}
}
