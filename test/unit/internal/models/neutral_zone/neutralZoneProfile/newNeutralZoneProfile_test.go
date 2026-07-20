package neutralZoneProfile_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
)

func TestWhenQualityVaries_SelectsMatchingProfileGuardMultiplier(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		subtestName        string
		quality            neutral_zone.Quality
		expectedMultiplier float64
	}{
		{"WhenQualityIsLowest_UsesLowestProfileGuardMultiplier", neutral_zone.QualityLowest, 0.9},
		{"WhenQualityIsLow_UsesLowProfileGuardMultiplier", neutral_zone.QualityLow, 1.1},
		{"WhenQualityIsMedium_UsesMediumProfileGuardMultiplier", neutral_zone.QualityMedium, 1.4},
		{"WhenQualityIsHigh_UsesHighProfileGuardMultiplier", neutral_zone.QualityHigh, 1.8},
		{"WhenQualityIsHighest_UsesHighestProfileGuardMultiplier", neutral_zone.QualityHighest, 2.3},
	}
	for _, testCase := range testCases {
		t.Run(testCase.subtestName, func(t *testing.T) {
			t.Parallel()
			// Arrange
			quality := testCase.quality

			// Act
			profile := neutral_zone.NewNeutralZoneProfile(quality)

			// Assert
			assert.InDelta(t, testCase.expectedMultiplier, profile.GuardMultiplier, test_helpers.Delta)
		})
	}
}

func TestWhenQualityVaries_SelectsMatchingCityGuardValues(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		subtestName          string
		quality              neutral_zone.Quality
		expectedPrimaryGuard int
	}{
		{"WhenQualityIsLowest_UsesLowestPrimaryCityGuard", neutral_zone.QualityLowest, 2000},
		{"WhenQualityIsLow_UsesLowPrimaryCityGuard", neutral_zone.QualityLow, 4000},
		{"WhenQualityIsMedium_UsesMediumPrimaryCityGuard", neutral_zone.QualityMedium, 8000},
		{"WhenQualityIsHigh_UsesHighPrimaryCityGuard", neutral_zone.QualityHigh, 16000},
		{"WhenQualityIsHighest_UsesHighestPrimaryCityGuard", neutral_zone.QualityHighest, 32000},
	}
	for _, testCase := range testCases {
		t.Run(testCase.subtestName, func(t *testing.T) {
			t.Parallel()
			// Arrange
			quality := testCase.quality

			// Act
			profile := neutral_zone.NewNeutralZoneProfile(quality)

			// Assert
			assert.Equal(t, testCase.expectedPrimaryGuard, profile.PrimaryCityGuardValue)
		})
	}
}

func TestWhenQualityIsLowest_UsesTier1GuardedPool(t *testing.T) {
	t.Parallel()
	// Arrange
	quality := neutral_zone.QualityLowest

	// Act
	profile := neutral_zone.NewNeutralZoneProfile(quality)

	// Assert
	assert.Equal(t, registry.GetGuardedContentPoolT1List(), profile.GuardedContentPool)
}

func TestWhenQualityIsLowest_UsesVeryPoorResourcesPool(t *testing.T) {
	t.Parallel()
	// Arrange
	quality := neutral_zone.QualityLowest

	// Act
	profile := neutral_zone.NewNeutralZoneProfile(quality)

	// Assert
	assert.Equal(t, []string{registry.GetResourcesContentPoolValues().StartZoneVeryPoor}, profile.ResourcesContentPool)
}

func TestWhenQualityIsLowest_UsesExtraPoorConstructionSids(t *testing.T) {
	t.Parallel()
	// Arrange
	quality := neutral_zone.QualityLowest

	// Act
	profile := neutral_zone.NewNeutralZoneProfile(quality)

	// Assert
	expected := registry.GetBuildingsConstructionSidValues().ExtraPoor
	assert.Equal(t, [2]string{expected, expected}, [2]string{profile.PrimaryBuildingsSid, profile.ExtraBuildingsSid})
}

func TestWhenQualityIsHighest_UsesDoubledTier5GuardedPool(t *testing.T) {
	t.Parallel()
	// Arrange
	quality := neutral_zone.QualityHighest

	// Act
	profile := neutral_zone.NewNeutralZoneProfile(quality)

	// Assert
	tier5List := registry.GetGuardedContentPoolT5List()
	assert.Equal(t, append(tier5List, tier5List...), profile.GuardedContentPool)
}

func TestWhenQualityIsHighest_UsesDoubledTier5UnguardedPool(t *testing.T) {
	t.Parallel()
	// Arrange
	quality := neutral_zone.QualityHighest

	// Act
	profile := neutral_zone.NewNeutralZoneProfile(quality)

	// Assert
	tier5List := registry.GetUnguardedContentPoolT5List()
	assert.Equal(t, append(tier5List, tier5List...), profile.UnguardedContentPool)
}

func TestWhenQualityIsHighest_UsesRichTreasureResourcesPool(t *testing.T) {
	t.Parallel()
	// Arrange
	quality := neutral_zone.QualityHighest

	// Act
	profile := neutral_zone.NewNeutralZoneProfile(quality)

	// Assert
	assert.Equal(t, []string{registry.GetResourcesContentPoolValues().TreasureZoneRich}, profile.ResourcesContentPool)
}

func TestWhenQualityIsHighest_UsesUltraRichConstructionSids(t *testing.T) {
	t.Parallel()
	// Arrange
	quality := neutral_zone.QualityHighest

	// Act
	profile := neutral_zone.NewNeutralZoneProfile(quality)

	// Assert
	expected := registry.GetBuildingsConstructionSidValues().UltraRich
	assert.Equal(t, [2]string{expected, expected}, [2]string{profile.PrimaryBuildingsSid, profile.ExtraBuildingsSid})
}

func TestWhenQualityIsHighest_UsesCenterLayout(t *testing.T) {
	t.Parallel()
	// Arrange
	quality := neutral_zone.QualityHighest

	// Act
	profile := neutral_zone.NewNeutralZoneProfile(quality)

	// Assert
	assert.Equal(t, registry.GetLayoutValues().Center, profile.Layout)
}
