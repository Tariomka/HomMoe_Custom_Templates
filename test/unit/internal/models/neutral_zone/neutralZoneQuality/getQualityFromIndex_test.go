package neutralZoneQuality_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/stretchr/testify/assert"
)

func TestWhenIndexVaries_ReturnsMatchingQuality(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		subtestName     string
		index           int
		expectedQuality neutral_zone.Quality
	}{
		{"WhenIndexIsZero_ReturnsLowestQuality", 0, neutral_zone.QualityLowest},
		{"WhenIndexIsOne_ReturnsLowQuality", 1, neutral_zone.QualityLow},
		{"WhenIndexIsTwo_ReturnsMediumQuality", 2, neutral_zone.QualityMedium},
		{"WhenIndexIsThree_ReturnsHighQuality", 3, neutral_zone.QualityHigh},
		{"WhenIndexIsFour_ReturnsHighestQuality", 4, neutral_zone.QualityHighest},
		{"WhenIndexIsNegative_ReturnsUnknownQuality", -1, neutral_zone.QualityUnknown},
		{"WhenIndexIsOutOfRange_ReturnsUnknownQuality", 5, neutral_zone.QualityUnknown},
	}
	for _, testCase := range testCases {
		t.Run(testCase.subtestName, func(t *testing.T) {
			t.Parallel()
			// Arrange
			index := testCase.index

			// Act
			quality := neutral_zone.GetQualityFromIndex(index)

			// Assert
			assert.Equal(t, testCase.expectedQuality, quality)
		})
	}
}
