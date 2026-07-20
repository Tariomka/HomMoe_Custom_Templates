package neutralZonePlan_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/stretchr/testify/assert"
)

func TestWhenQualityAndCastleCountVary_ComputesScoreFromQualityAndCappedCastles(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		subtestName   string
		plan          neutral_zone.Plan
		expectedScore float64
	}{
		{
			"WhenQualityIsLowestWithoutCastles_ScoresHalf",
			neutral_zone.Plan{Quality: neutral_zone.QualityLowest, CastleCount: 0},
			0.5,
		},
		{
			"WhenQualityIsLowWithoutCastles_ScoresOne",
			neutral_zone.Plan{Quality: neutral_zone.QualityLow, CastleCount: 0},
			1.0,
		},
		{
			"WhenQualityIsMediumWithoutCastles_ScoresTwo",
			neutral_zone.Plan{Quality: neutral_zone.QualityMedium, CastleCount: 0},
			2.0,
		},
		{
			"WhenQualityIsHighWithoutCastles_ScoresThree",
			neutral_zone.Plan{Quality: neutral_zone.QualityHigh, CastleCount: 0},
			3.0,
		},
		{
			"WhenQualityIsHighestWithoutCastles_ScoresFour",
			neutral_zone.Plan{Quality: neutral_zone.QualityHighest, CastleCount: 0},
			4.0,
		},
		{
			"WhenLowZoneHasTwoCastles_AddsFifteenHundredthsPerCastle",
			neutral_zone.Plan{Quality: neutral_zone.QualityLow, CastleCount: 2},
			1.3,
		},
		{
			"WhenCastleCountExceedsFour_CapsCastleBonusAtFour",
			neutral_zone.Plan{Quality: neutral_zone.QualityHigh, CastleCount: 9},
			3.6,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.subtestName, func(t *testing.T) {
			t.Parallel()
			// Arrange
			plan := testCase.plan

			// Act
			score := plan.GetBalanceScore()

			// Assert
			assert.InDelta(t, testCase.expectedScore, score, 1e-9)
		})
	}
}
