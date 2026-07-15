package neutralZonePlan_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutralZone"
	"github.com/stretchr/testify/assert"
)

func TestWhenQualityAndCastleCountVary_ComputesScoreFromQualityAndCappedCastles(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		subtestName   string
		plan          neutralZone.Plan
		expectedScore float64
	}{
		{
			"WhenQualityIsLowestWithoutCastles_ScoresHalf",
			neutralZone.Plan{Quality: neutralZone.QualityLowest, CastleCount: 0},
			0.5,
		},
		{
			"WhenQualityIsLowWithoutCastles_ScoresOne",
			neutralZone.Plan{Quality: neutralZone.QualityLow, CastleCount: 0},
			1.0,
		},
		{
			"WhenQualityIsMediumWithoutCastles_ScoresTwo",
			neutralZone.Plan{Quality: neutralZone.QualityMedium, CastleCount: 0},
			2.0,
		},
		{
			"WhenQualityIsHighWithoutCastles_ScoresThree",
			neutralZone.Plan{Quality: neutralZone.QualityHigh, CastleCount: 0},
			3.0,
		},
		{
			"WhenQualityIsHighestWithoutCastles_ScoresFour",
			neutralZone.Plan{Quality: neutralZone.QualityHighest, CastleCount: 0},
			4.0,
		},
		{
			"WhenLowZoneHasTwoCastles_AddsFifteenHundredthsPerCastle",
			neutralZone.Plan{Quality: neutralZone.QualityLow, CastleCount: 2},
			1.3,
		},
		{
			"WhenCastleCountExceedsFour_CapsCastleBonusAtFour",
			neutralZone.Plan{Quality: neutralZone.QualityHigh, CastleCount: 9},
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
