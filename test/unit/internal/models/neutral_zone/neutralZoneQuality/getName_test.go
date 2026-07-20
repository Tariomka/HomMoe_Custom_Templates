package neutralZoneQuality_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/stretchr/testify/assert"
)

func TestWhenQualityVaries_ReturnsMatchingName(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		subtestName  string
		quality      neutral_zone.Quality
		expectedName string
	}{
		{"WhenQualityIsHighest_ReturnsPlatinum", neutral_zone.QualityHighest, "Platinum"},
		{"WhenQualityIsHigh_ReturnsGold", neutral_zone.QualityHigh, "Gold"},
		{"WhenQualityIsMedium_ReturnsSilver", neutral_zone.QualityMedium, "Silver"},
		{"WhenQualityIsLow_ReturnsBronze", neutral_zone.QualityLow, "Bronze"},
		{"WhenQualityIsLowest_ReturnsPlastic", neutral_zone.QualityLowest, "Plastic"},
		{"WhenQualityIsUnknown_ReturnsUnknown", neutral_zone.QualityUnknown, "Unknown"},
	}
	for _, testCase := range testCases {
		t.Run(testCase.subtestName, func(t *testing.T) {
			t.Parallel()
			// Arrange
			quality := testCase.quality

			// Act
			name := quality.GetName()

			// Assert
			assert.Equal(t, testCase.expectedName, name)
		})
	}
}
