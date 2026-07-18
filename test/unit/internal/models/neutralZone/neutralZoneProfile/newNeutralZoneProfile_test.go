package neutralZoneProfile_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutralZone"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
)

func TestWhenQualityVaries_SelectsMatchingProfileGuardMultiplier(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		subtestName        string
		quality            neutralZone.Quality
		expectedMultiplier float64
	}{
		{"WhenQualityIsLowest_UsesLowestProfileGuardMultiplier", neutralZone.QualityLowest, 0.9},
		{"WhenQualityIsLow_UsesLowProfileGuardMultiplier", neutralZone.QualityLow, 1.1},
		{"WhenQualityIsMedium_UsesMediumProfileGuardMultiplier", neutralZone.QualityMedium, 1.4},
		{"WhenQualityIsHigh_UsesHighProfileGuardMultiplier", neutralZone.QualityHigh, 1.8},
		{"WhenQualityIsHighest_UsesHighestProfileGuardMultiplier", neutralZone.QualityHighest, 2.3},
	}
	for _, testCase := range testCases {
		t.Run(testCase.subtestName, func(t *testing.T) {
			t.Parallel()
			// Arrange
			quality := testCase.quality

			// Act
			profile := neutralZone.NewNeutralZoneProfile(quality)

			// Assert
			assert.InDelta(t, testCase.expectedMultiplier, profile.GuardMultiplier, test_helpers.Delta)
		})
	}
}

func TestWhenQualityVaries_SelectsMatchingCityGuardValues(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		subtestName          string
		quality              neutralZone.Quality
		expectedPrimaryGuard int
	}{
		{"WhenQualityIsLowest_UsesLowestPrimaryCityGuard", neutralZone.QualityLowest, 2000},
		{"WhenQualityIsLow_UsesLowPrimaryCityGuard", neutralZone.QualityLow, 4000},
		{"WhenQualityIsMedium_UsesMediumPrimaryCityGuard", neutralZone.QualityMedium, 8000},
		{"WhenQualityIsHigh_UsesHighPrimaryCityGuard", neutralZone.QualityHigh, 16000},
		{"WhenQualityIsHighest_UsesHighestPrimaryCityGuard", neutralZone.QualityHighest, 32000},
	}
	for _, testCase := range testCases {
		t.Run(testCase.subtestName, func(t *testing.T) {
			t.Parallel()
			// Arrange
			quality := testCase.quality

			// Act
			profile := neutralZone.NewNeutralZoneProfile(quality)

			// Assert
			assert.Equal(t, testCase.expectedPrimaryGuard, profile.PrimaryCityGuardValue)
		})
	}
}

func TestWhenQualityIsLowest_UsesTier1GuardedPool(t *testing.T) {
	t.Parallel()
	// Arrange
	quality := neutralZone.QualityLowest

	// Act
	profile := neutralZone.NewNeutralZoneProfile(quality)

	// Assert
	assert.Equal(t, registry.GetGuardedContentPoolT1List(), profile.GuardedContentPool)
}

func TestWhenQualityIsLowest_UsesVeryPoorResourcesPool(t *testing.T) {
	t.Parallel()
	// Arrange
	quality := neutralZone.QualityLowest

	// Act
	profile := neutralZone.NewNeutralZoneProfile(quality)

	// Assert
	assert.Equal(t, []string{registry.GetResourcesContentPoolValues().StartZoneVeryPoor}, profile.ResourcesContentPool)
}

func TestWhenQualityIsLowest_UsesExtraPoorConstructionSids(t *testing.T) {
	t.Parallel()
	// Arrange
	quality := neutralZone.QualityLowest

	// Act
	profile := neutralZone.NewNeutralZoneProfile(quality)

	// Assert
	expected := registry.GetBuildingsConstructionSidValues().ExtraPoor
	assert.Equal(t, [2]string{expected, expected}, [2]string{profile.PrimaryBuildingsSid, profile.ExtraBuildingsSid})
}

func TestWhenQualityIsHighest_UsesDoubledTier5GuardedPool(t *testing.T) {
	t.Parallel()
	// Arrange
	quality := neutralZone.QualityHighest

	// Act
	profile := neutralZone.NewNeutralZoneProfile(quality)

	// Assert
	tier5List := registry.GetGuardedContentPoolT5List()
	assert.Equal(t, append(tier5List, tier5List...), profile.GuardedContentPool)
}

func TestWhenQualityIsHighest_UsesDoubledTier5UnguardedPool(t *testing.T) {
	t.Parallel()
	// Arrange
	quality := neutralZone.QualityHighest

	// Act
	profile := neutralZone.NewNeutralZoneProfile(quality)

	// Assert
	tier5List := registry.GetUnguardedContentPoolT5List()
	assert.Equal(t, append(tier5List, tier5List...), profile.UnguardedContentPool)
}

func TestWhenQualityIsHighest_UsesRichTreasureResourcesPool(t *testing.T) {
	t.Parallel()
	// Arrange
	quality := neutralZone.QualityHighest

	// Act
	profile := neutralZone.NewNeutralZoneProfile(quality)

	// Assert
	assert.Equal(t, []string{registry.GetResourcesContentPoolValues().TreasureZoneRich}, profile.ResourcesContentPool)
}

func TestWhenQualityIsHighest_UsesUltraRichConstructionSids(t *testing.T) {
	t.Parallel()
	// Arrange
	quality := neutralZone.QualityHighest

	// Act
	profile := neutralZone.NewNeutralZoneProfile(quality)

	// Assert
	expected := registry.GetBuildingsConstructionSidValues().UltraRich
	assert.Equal(t, [2]string{expected, expected}, [2]string{profile.PrimaryBuildingsSid, profile.ExtraBuildingsSid})
}

func TestWhenQualityIsHighest_UsesCenterLayout(t *testing.T) {
	t.Parallel()
	// Arrange
	quality := neutralZone.QualityHighest

	// Act
	profile := neutralZone.NewNeutralZoneProfile(quality)

	// Assert
	assert.Equal(t, registry.GetLayoutValues().Center, profile.Layout)
}
