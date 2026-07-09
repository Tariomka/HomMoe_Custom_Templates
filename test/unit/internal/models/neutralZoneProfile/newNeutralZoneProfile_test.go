package neutralZoneProfile_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/test/helpers"
	"github.com/stretchr/testify/assert"
)

func TestWhenQualityVaries_SelectsMatchingProfileGuardMultiplier(t *testing.T) {
	testCases := []struct {
		subtestName        string
		quality            models.NeutralZoneQuality
		expectedMultiplier float64
	}{
		{"WhenQualityIsLow_UsesLowProfileGuardMultiplier", models.QualityLow, 1.1},
		{"WhenQualityIsMedium_UsesMediumProfileGuardMultiplier", models.QualityMedium, 1.4},
		{"WhenQualityIsHigh_UsesHighProfileGuardMultiplier", models.QualityHigh, 1.8},
	}
	for _, testCase := range testCases {
		t.Run(testCase.subtestName, func(t *testing.T) {
			// Arrange
			quality := testCase.quality

			// Act
			profile := models.NewNeutralZoneProfile(quality)

			// Assert
			assert.InDelta(t, testCase.expectedMultiplier, profile.GuardMultiplier, helpers.Delta)
		})
	}
}

func TestWhenQualityVaries_SelectsMatchingCityGuardValues(t *testing.T) {
	testCases := []struct {
		subtestName          string
		quality              models.NeutralZoneQuality
		expectedPrimaryGuard int
	}{
		{"WhenQualityIsLow_UsesLowPrimaryCityGuard", models.QualityLow, 4000},
		{"WhenQualityIsMedium_UsesMediumPrimaryCityGuard", models.QualityMedium, 8000},
		{"WhenQualityIsHigh_UsesHighPrimaryCityGuard", models.QualityHigh, 16000},
	}
	for _, testCase := range testCases {
		t.Run(testCase.subtestName, func(t *testing.T) {
			// Arrange
			quality := testCase.quality

			// Act
			profile := models.NewNeutralZoneProfile(quality)

			// Assert
			assert.Equal(t, testCase.expectedPrimaryGuard, profile.PrimaryCityGuardValue)
		})
	}
}
