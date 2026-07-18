package neutralZoneQuality_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutralZone"
	"github.com/stretchr/testify/assert"
)

func TestWhenQualityVaries_ReturnsMatchingBalanceScore(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		subtestName   string
		quality       neutralZone.Quality
		expectedScore float64
	}{
		{"WhenQualityIsHighest_ReturnsFour", neutralZone.QualityHighest, 4.0},
		{"WhenQualityIsHigh_ReturnsThree", neutralZone.QualityHigh, 3.0},
		{"WhenQualityIsMedium_ReturnsTwo", neutralZone.QualityMedium, 2.0},
		{"WhenQualityIsLow_ReturnsOne", neutralZone.QualityLow, 1.0},
		{"WhenQualityIsLowest_ReturnsHalf", neutralZone.QualityLowest, 0.5},
		{"WhenQualityIsUnknown_ReturnsZero", neutralZone.QualityUnknown, 0.0},
	}
	for _, testCase := range testCases {
		t.Run(testCase.subtestName, func(t *testing.T) {
			t.Parallel()
			// Arrange
			quality := testCase.quality

			// Act
			score := quality.GetBalanceScore()

			// Assert
			assert.InDelta(t, testCase.expectedScore, score, 0)
		})
	}
}
