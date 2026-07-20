package neutralZonePlan_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/stretchr/testify/assert"
)

func TestWhenLabelQualityVaries_MapsQualityToTier(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		subtestName  string
		quality      neutral_zone.Quality
		expectedTier int
	}{
		{"WhenLabelHasHighestQuality_ReturnsTierFour", neutral_zone.QualityHighest, 4},
		{"WhenLabelHasHighQuality_ReturnsTierThree", neutral_zone.QualityHigh, 3},
		{"WhenLabelHasMediumQuality_ReturnsTierTwo", neutral_zone.QualityMedium, 2},
		{"WhenLabelHasLowQuality_ReturnsTierOne", neutral_zone.QualityLow, 1},
		{"WhenLabelHasLowestQuality_ReturnsTierZero", neutral_zone.QualityLowest, 0},
	}
	for _, testCase := range testCases {
		t.Run(testCase.subtestName, func(t *testing.T) {
			t.Parallel()
			// Arrange
			plans := neutral_zone.Plans{{Label: "A", Quality: testCase.quality}}

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
	plans := neutral_zone.Plans{{Label: "A", Quality: neutral_zone.QualityHigh}}

	// Act
	tier := plans.GetTier("missing")

	// Assert
	assert.Equal(t, neutral_zone.QualityUnknown.GetIndex(), tier)
}
