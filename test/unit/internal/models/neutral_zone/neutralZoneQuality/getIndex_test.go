package neutralZoneQuality_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/stretchr/testify/assert"
)

func TestWhenQualityVaries_ReturnsMatchingIndex(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		subtestName   string
		quality       neutral_zone.Quality
		expectedIndex int
	}{
		{"WhenQualityIsUnknown_ReturnsMinusOne", neutral_zone.QualityUnknown, -1},
		{"WhenQualityIsLowest_ReturnsZero", neutral_zone.QualityLowest, 0},
		{"WhenQualityIsLow_ReturnsOne", neutral_zone.QualityLow, 1},
		{"WhenQualityIsMedium_ReturnsTwo", neutral_zone.QualityMedium, 2},
		{"WhenQualityIsHigh_ReturnsThree", neutral_zone.QualityHigh, 3},
		{"WhenQualityIsHighest_ReturnsFour", neutral_zone.QualityHighest, 4},
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
