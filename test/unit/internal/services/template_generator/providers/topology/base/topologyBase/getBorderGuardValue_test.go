package topologyBase_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/base"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
)

func TestWhenBothLabelsArePlayers_GuardValueIsPlayerBorderStrength(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())

	// Act
	guardValue := topologyBase.GetBorderGuardValue(
		"A", "B", []string{"A", "B"}, nil, newUnitTuning())

	// Assert
	assert.Equal(t, 30000, guardValue)
}

func TestWhenBothLabelsAreNeutral_HigherQualityGuardWins(t *testing.T) {
	t.Parallel()
	// Arrange
	neutralPlans := neutral_zone.Plans{
		{Label: "C", Quality: neutral_zone.QualityLow, CastleCount: 0},
		{Label: "D", Quality: neutral_zone.QualityHigh, CastleCount: 0},
	}
	testCases := []struct {
		name       string
		firstLabel string
		otherLabel string
	}{
		{name: "WhenSecondLabelHasHigherQuality_ItsGuardIsUsed", firstLabel: "C", otherLabel: "D"},
		{name: "WhenFirstLabelHasHigherQuality_ItsGuardIsUsed", firstLabel: "D", otherLabel: "C"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			// Arrange
			topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())

			// Act
			guardValue := topologyBase.GetBorderGuardValue(
				testCase.firstLabel, testCase.otherLabel, []string{"A", "B"}, neutralPlans, newUnitTuning())

			// Assert
			assert.Equal(t, 25000, guardValue)
		})
	}
}

func TestWhenOnlyFirstLabelIsPlayer_NeutralSecondLabelQualityDrivesGuard(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())
	neutralPlans := neutral_zone.Plans{
		{Label: "C", Quality: neutral_zone.QualityLow, CastleCount: 0},
	}

	// Act
	guardValue := topologyBase.GetBorderGuardValue(
		"A", "C", []string{"A", "B"}, neutralPlans, newUnitTuning())

	// Assert
	assert.Equal(t, 15000, guardValue)
}

func TestWhenOnlySecondLabelIsPlayer_NeutralFirstLabelQualityDrivesGuard(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())
	neutralPlans := neutral_zone.Plans{
		{Label: "C", Quality: neutral_zone.QualityLow, CastleCount: 0},
	}

	// Act
	guardValue := topologyBase.GetBorderGuardValue(
		"C", "A", []string{"A", "B"}, neutralPlans, newUnitTuning())

	// Assert
	assert.Equal(t, 15000, guardValue)
}

func TestWhenNeutralLabelHasNoPlan_UsesUnknownQualityGuard(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())

	// Act
	guardValue := topologyBase.GetBorderGuardValue(
		"A", "Z", []string{"A", "B"}, nil, newUnitTuning())

	// Assert
	assert.Equal(t, 30000, guardValue)
}

func TestWhenBorderGuardMultiplierIsDoubled_PlayerBorderGuardIsScaled(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())
	tuning := newUnitTuning()
	tuning.BorderGuardStrengthMultiplier = 2.0

	// Act
	guardValue := topologyBase.GetBorderGuardValue(
		"A", "B", []string{"A", "B"}, nil, tuning)

	// Assert
	assert.Equal(t, 60000, guardValue)
}
