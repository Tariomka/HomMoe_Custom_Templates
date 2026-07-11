package neutralZonePlan_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestWhenLabelQualityVaries_MapsQualityToTier(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		subtestName  string
		quality      models.NeutralZoneQuality
		expectedTier int
	}{
		{"WhenLabelHasHighQuality_ReturnsTierThree", models.QualityHigh, 3},
		{"WhenLabelHasMediumQuality_ReturnsTierTwo", models.QualityMedium, 2},
		{"WhenLabelHasLowQuality_ReturnsTierOne", models.QualityLow, 1},
	}
	for _, testCase := range testCases {
		t.Run(testCase.subtestName, func(t *testing.T) {
			t.Parallel()
			// Arrange
			plans := models.NeutralZonePlans{{Label: "A", Quality: testCase.quality}}

			// Act
			tier := plans.GetTier("A")

			// Assert
			assert.Equal(t, testCase.expectedTier, tier)
		})
	}
}

func TestWhenLabelIsNotFound_ReturnsTierOne(t *testing.T) {
	t.Parallel()
	// Arrange
	plans := models.NeutralZonePlans{{Label: "A", Quality: models.QualityHigh}}

	// Act
	tier := plans.GetTier("missing")

	// Assert
	assert.Equal(t, 1, tier)
}
