package neutralZoneQuality_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutralZone"
	"github.com/stretchr/testify/assert"
)

func TestWhenQualityVaries_ReturnsMatchingIndex(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		subtestName   string
		quality       neutralZone.Quality
		expectedIndex int
	}{
		{"WhenQualityIsUnknown_ReturnsMinusOne", neutralZone.QualityUnknown, -1},
		{"WhenQualityIsLowest_ReturnsZero", neutralZone.QualityLowest, 0},
		{"WhenQualityIsLow_ReturnsOne", neutralZone.QualityLow, 1},
		{"WhenQualityIsMedium_ReturnsTwo", neutralZone.QualityMedium, 2},
		{"WhenQualityIsHigh_ReturnsThree", neutralZone.QualityHigh, 3},
		{"WhenQualityIsHighest_ReturnsFour", neutralZone.QualityHighest, 4},
	}
	for _, testCase := range testCases {
		t.Run(testCase.subtestName, func(t *testing.T) {
			t.Parallel()
			// Arrange
			quality := testCase.quality

			// Act
			index := quality.GetIndex()

			// Assert
			assert.Equal(t, testCase.expectedIndex, index)
		})
	}
}
