package neutralZonePlan_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestWhenQualityAndCastleCountVary_ComputesScoreFromQualityAndCappedCastles(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		subtestName   string
		plan          models.NeutralZonePlan
		expectedScore float64
	}{
		{
			"WhenQualityIsLowWithoutCastles_ScoresOne",
			models.NeutralZonePlan{Quality: models.QualityLow, CastleCount: 0},
			1.0,
		},
		{
			"WhenQualityIsMediumWithoutCastles_ScoresTwo",
			models.NeutralZonePlan{Quality: models.QualityMedium, CastleCount: 0},
			2.0,
		},
		{
			"WhenQualityIsHighWithoutCastles_ScoresThree",
			models.NeutralZonePlan{Quality: models.QualityHigh, CastleCount: 0},
			3.0,
		},
		{
			"WhenLowZoneHasTwoCastles_AddsFifteenHundredthsPerCastle",
			models.NeutralZonePlan{Quality: models.QualityLow, CastleCount: 2},
			1.3,
		},
		{
			"WhenCastleCountExceedsFour_CapsCastleBonusAtFour",
			models.NeutralZonePlan{Quality: models.QualityHigh, CastleCount: 9},
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
