package neutralZonePlan_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutralZone"
	"github.com/stretchr/testify/assert"
)

func TestWhenLabelQualityVaries_MapsQualityToTier(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		subtestName  string
		quality      neutralZone.Quality
		expectedTier int
	}{
		{"WhenLabelHasHighestQuality_ReturnsTierFour", neutralZone.QualityHighest, 4},
		{"WhenLabelHasHighQuality_ReturnsTierThree", neutralZone.QualityHigh, 3},
		{"WhenLabelHasMediumQuality_ReturnsTierTwo", neutralZone.QualityMedium, 2},
		{"WhenLabelHasLowQuality_ReturnsTierOne", neutralZone.QualityLow, 1},
		{"WhenLabelHasLowestQuality_ReturnsTierZero", neutralZone.QualityLowest, 0},
	}
	for _, testCase := range testCases {
		t.Run(testCase.subtestName, func(t *testing.T) {
			t.Parallel()
			// Arrange
			plans := neutralZone.Plans{{Label: "A", Quality: testCase.quality}}

			// Act
			tier := plans.GetTier("A")

			// Assert
			assert.Equal(t, testCase.expectedTier, tier)
		})
	}
}

func TestWhenLabelIsNotFound_ReturnsUnknownTierIndex(t *testing.T) {
	t.Parallel()
	// Arrange
	plans := neutralZone.Plans{{Label: "A", Quality: neutralZone.QualityHigh}}

	// Act
	tier := plans.GetTier("missing")

	// Assert
	assert.Equal(t, neutralZone.QualityUnknown.GetIndex(), tier)
}
