package guardStrength_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_connections"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/stretchr/testify/assert"
)

func TestWhenQualityVaries_ReturnsMatchingGuardStrength(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		subtestName string
		quality     neutral_zone.Quality
		expected    common_connections.GuardStrength
	}{
		{
			"WhenQualityIsLowest_ReturnsBronzeStrength",
			neutral_zone.QualityLowest,
			common_connections.GetBronzeGuardStrength(),
		},
		{
			"WhenQualityIsLow_ReturnsBronzeStrength",
			neutral_zone.QualityLow,
			common_connections.GetBronzeGuardStrength(),
		},
		{
			"WhenQualityIsMedium_ReturnsSilverStrength",
			neutral_zone.QualityMedium,
			common_connections.GetSilverGuardStrength(),
		},
		{"WhenQualityIsHigh_ReturnsGoldStrength", neutral_zone.QualityHigh, common_connections.GetGoldGuardStrength()},
		{
			"WhenQualityIsHighest_ReturnsHubStrength",
			neutral_zone.QualityHighest,
			common_connections.GetHubGuardStrength(),
		},
		{
			"WhenQualityIsUnknown_ReturnsPlayerToPlayerStrength",
			neutral_zone.QualityUnknown,
			common_connections.GetPlayerToPlayerGuardStrength(),
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.subtestName, func(t *testing.T) {
			t.Parallel()
			// Arrange
			quality := testCase.quality

			// Act
			strength := common_connections.GetGuardStrengthForQuality(quality)

			// Assert
			assert.Equal(t, testCase.expected, strength)
		})
	}
}
