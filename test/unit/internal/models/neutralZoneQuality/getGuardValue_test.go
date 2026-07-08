package neutralZoneQuality_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestWhenQualityVaries_ReturnsMatchingGuardValue(t *testing.T) {
	testCases := []struct {
		subtestName        string
		quality            models.NeutralZoneQuality
		expectedGuardValue int
	}{
		{"WhenQualityIsHigh_ReturnsTwentyFiveThousand", models.QualityHigh, 25_000},
		{"WhenQualityIsMedium_ReturnsTwentyThousand", models.QualityMedium, 20_000},
		{"WhenQualityIsLow_ReturnsFifteenThousand", models.QualityLow, 15_000},
	}
	for _, testCase := range testCases {
		t.Run(testCase.subtestName, func(t *testing.T) {
			// Arrange
			quality := testCase.quality

			// Act
			guardValue := quality.GetGuardValue()

			// Assert
			assert.Equal(t, testCase.expectedGuardValue, guardValue)
		})
	}
}
