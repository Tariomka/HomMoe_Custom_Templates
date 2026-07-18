package zoneConfig_test

import (
	"math"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
)

func TestWhenRandomizationInputVaries_ReturnsEffectiveValue(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		subtestName string
		settings    config.ZoneConfig
		expected    float64
	}{
		{
			"WhenAdvancedModeIsDisabled_ReturnsDefaultRandomization",
			config.ZoneConfig{Advanced: config.AdvancedSettings{Enabled: false}, GuardRandomization: 0.4},
			0.05,
		},
		{
			"WhenValueIsNaN_ReturnsDefaultRandomization",
			config.ZoneConfig{Advanced: config.AdvancedSettings{Enabled: true}, GuardRandomization: math.NaN()},
			0.05,
		},
		{
			"WhenValueIsPositiveInfinity_ReturnsDefaultRandomization",
			config.ZoneConfig{Advanced: config.AdvancedSettings{Enabled: true}, GuardRandomization: math.Inf(1)},
			0.05,
		},
		{
			"WhenValueIsNegativeInfinity_ReturnsDefaultRandomization",
			config.ZoneConfig{Advanced: config.AdvancedSettings{Enabled: true}, GuardRandomization: math.Inf(-1)},
			0.05,
		},
		{
			"WhenValueIsNegative_ClampsToZero",
			config.ZoneConfig{Advanced: config.AdvancedSettings{Enabled: true}, GuardRandomization: -0.3},
			0,
		},
		{
			"WhenValueExceedsHalf_ClampsToHalf",
			config.ZoneConfig{Advanced: config.AdvancedSettings{Enabled: true}, GuardRandomization: 0.75},
			0.5,
		},
		{
			"WhenValueIsWithinRange_ReturnsValueUnchanged",
			config.ZoneConfig{Advanced: config.AdvancedSettings{Enabled: true}, GuardRandomization: 0.25},
			0.25,
		},
		{
			"WhenValueHasExtraPrecision_RoundsToThreeDecimals",
			config.ZoneConfig{Advanced: config.AdvancedSettings{Enabled: true}, GuardRandomization: 0.1234},
			0.123,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.subtestName, func(t *testing.T) {
			t.Parallel()
			// Arrange

			// Act
			actual := testCase.settings.GetEffectiveGuardRandomization()

			// Assert
			assert.InDelta(t, testCase.expected, actual, test_helpers.Delta)
		})
	}
}
