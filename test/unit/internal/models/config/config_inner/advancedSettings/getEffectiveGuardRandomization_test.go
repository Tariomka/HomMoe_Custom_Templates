package advancedSettings_test

import (
	"math"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/config/config_inner"
	"github.com/Tariomka/hommoe_custom_templates/test/helpers"
	"github.com/stretchr/testify/assert"
)

func TestWhenRandomizationInputVaries_ReturnsEffectiveValue(t *testing.T) {
	testCases := []struct {
		subtestName string
		settings    config_inner.AdvancedSettings
		expected    float64
	}{
		{
			"WhenAdvancedModeIsDisabled_ReturnsDefaultRandomization",
			config_inner.AdvancedSettings{Enabled: false, GuardRandomization: 0.4},
			0.05,
		},
		{
			"WhenValueIsNaN_ReturnsDefaultRandomization",
			config_inner.AdvancedSettings{Enabled: true, GuardRandomization: math.NaN()},
			0.05,
		},
		{
			"WhenValueIsPositiveInfinity_ReturnsDefaultRandomization",
			config_inner.AdvancedSettings{Enabled: true, GuardRandomization: math.Inf(1)},
			0.05,
		},
		{
			"WhenValueIsNegativeInfinity_ReturnsDefaultRandomization",
			config_inner.AdvancedSettings{Enabled: true, GuardRandomization: math.Inf(-1)},
			0.05,
		},
		{
			"WhenValueIsNegative_ClampsToZero",
			config_inner.AdvancedSettings{Enabled: true, GuardRandomization: -0.3},
			0,
		},
		{
			"WhenValueExceedsHalf_ClampsToHalf",
			config_inner.AdvancedSettings{Enabled: true, GuardRandomization: 0.75},
			0.5,
		},
		{
			"WhenValueIsWithinRange_ReturnsValueUnchanged",
			config_inner.AdvancedSettings{Enabled: true, GuardRandomization: 0.25},
			0.25,
		},
		{
			"WhenValueHasExtraPrecision_RoundsToThreeDecimals",
			config_inner.AdvancedSettings{Enabled: true, GuardRandomization: 0.1234},
			0.123,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.subtestName, func(t *testing.T) {
			// Arrange - settings provided by the test case.

			// Act
			actual := testCase.settings.GetEffectiveGuardRandomization()

			// Assert
			assert.InDelta(t, testCase.expected, actual, helpers.Delta)
		})
	}
}
