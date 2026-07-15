package neutralZoneQuality_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutralZone"
	"github.com/stretchr/testify/assert"
)

func TestWhenQualityVaries_ReturnsMatchingGuardValue(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		subtestName        string
		quality            neutralZone.Quality
		expectedGuardValue int
	}{
		{"WhenQualityIsHighest_ReturnsThirtyThousand", neutralZone.QualityHighest, 30_000},
		{"WhenQualityIsHigh_ReturnsTwentyFiveThousand", neutralZone.QualityHigh, 25_000},
		{"WhenQualityIsMedium_ReturnsTwentyThousand", neutralZone.QualityMedium, 20_000},
		{"WhenQualityIsLow_ReturnsFifteenThousand", neutralZone.QualityLow, 15_000},
		{"WhenQualityIsLowest_ReturnsTenThousand", neutralZone.QualityLowest, 10_000},
	}
	for _, testCase := range testCases {
		t.Run(testCase.subtestName, func(t *testing.T) {
			t.Parallel()
			// Arrange
			quality := testCase.quality

			// Act
			guardValue := quality.GetGuardValue()

			// Assert
			assert.Equal(t, testCase.expectedGuardValue, guardValue)
		})
	}
}
