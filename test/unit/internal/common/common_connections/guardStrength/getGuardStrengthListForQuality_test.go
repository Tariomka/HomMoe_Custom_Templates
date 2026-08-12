package guardStrength_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_connections"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/stretchr/testify/assert"
)

func newExpectedStrengthList(strength models.GuardStrength) []data.Tuple[string, int] {
	return []data.Tuple[string, int]{
		data.NewTuple("Default", strength.Default),
		data.NewTuple("Weakest", strength.Weakest),
		data.NewTuple("Low", strength.Low),
		data.NewTuple("Medium", strength.Medium),
		data.NewTuple("High", strength.High),
		data.NewTuple("Very High", strength.VeryHigh),
	}
}

func TestWhenQualityVaries_ReturnsMatchingGuardStrengthList(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		subtestName string
		quality     neutral_zone.Quality
		expected    []data.Tuple[string, int]
	}{
		{
			"WhenQualityIsLowest_ReturnsBronzeStrengthList",
			neutral_zone.QualityLowest,
			newExpectedStrengthList(common_connections.GetBronzeGuardStrength()),
		},
		{
			"WhenQualityIsLow_ReturnsBronzeStrengthList",
			neutral_zone.QualityLow,
			newExpectedStrengthList(common_connections.GetBronzeGuardStrength()),
		},
		{
			"WhenQualityIsMedium_ReturnsSilverStrengthList",
			neutral_zone.QualityMedium,
			newExpectedStrengthList(common_connections.GetSilverGuardStrength()),
		},
		{
			"WhenQualityIsHigh_ReturnsGoldStrengthList",
			neutral_zone.QualityHigh,
			newExpectedStrengthList(common_connections.GetGoldGuardStrength()),
		},
		{
			"WhenQualityIsHighest_ReturnsHubStrengthList",
			neutral_zone.QualityHighest,
			newExpectedStrengthList(common_connections.GetHubGuardStrength()),
		},
		{
			"WhenQualityIsUnknown_ReturnsPlayerToPlayerStrengthList",
			neutral_zone.QualityUnknown,
			newExpectedStrengthList(common_connections.GetPlayerToPlayerGuardStrength()),
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.subtestName, func(t *testing.T) {
			t.Parallel()
			// Arrange
			quality := testCase.quality

			// Act
			strengthList := common_connections.GetGuardStrengthListForQuality(quality)

			// Assert
			assert.Equal(t, testCase.expected, strengthList)
		})
	}
}
