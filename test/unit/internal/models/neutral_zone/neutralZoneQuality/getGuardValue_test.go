package neutralZoneQuality_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/stretchr/testify/assert"
)

func TestWhenQualityVaries_ReturnsMatchingGuardValue(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		subtestName        string
		quality            neutral_zone.Quality
		expectedGuardValue int
	}{
		{"WhenQualityIsHighest_ReturnsThirtyFiveThousand", neutral_zone.QualityHighest, 35_000},
		{"WhenQualityIsHigh_ReturnsTwentyFiveThousand", neutral_zone.QualityHigh, 25_000},
		{"WhenQualityIsMedium_ReturnsTwentyThousand", neutral_zone.QualityMedium, 20_000},
		{"WhenQualityIsLow_ReturnsFifteenThousand", neutral_zone.QualityLow, 15_000},
		{"WhenQualityIsLowest_ReturnsTenThousand", neutral_zone.QualityLowest, 10_000},
		{"WhenQualityIsUnknown_ReturnsThirtyThousand", neutral_zone.QualityUnknown, 30_000},
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
