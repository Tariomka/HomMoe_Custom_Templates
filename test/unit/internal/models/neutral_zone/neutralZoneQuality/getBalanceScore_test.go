package neutralZoneQuality_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/stretchr/testify/assert"
)

func TestWhenQualityVaries_ReturnsMatchingBalanceScore(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		subtestName   string
		quality       neutral_zone.Quality
		expectedScore float64
	}{
		{"WhenQualityIsHighest_ReturnsFour", neutral_zone.QualityHighest, 4.0},
		{"WhenQualityIsHigh_ReturnsThree", neutral_zone.QualityHigh, 3.0},
		{"WhenQualityIsMedium_ReturnsTwo", neutral_zone.QualityMedium, 2.0},
		{"WhenQualityIsLow_ReturnsOne", neutral_zone.QualityLow, 1.0},
		{"WhenQualityIsLowest_ReturnsHalf", neutral_zone.QualityLowest, 0.5},
		{"WhenQualityIsUnknown_ReturnsZero", neutral_zone.QualityUnknown, 0.0},
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
