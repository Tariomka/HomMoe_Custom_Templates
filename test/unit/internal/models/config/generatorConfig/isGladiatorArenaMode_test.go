package generatorConfig_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/stretchr/testify/assert"
)

func TestWhenGladiatorRuleAndVictoryConditionCombinationsVary_ReportsArenaModeAccordingly(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		subtestName string
		mutate      func(configuration *config.GeneratorConfig)
		expected    bool
	}{
		{
			"WhenGladiatorArenaRulesAreEnabled_ReturnsTrue",
			func(configuration *config.GeneratorConfig) {
				configuration.GladiatorArenaRules.Enabled = true
			},
			true,
		},
		{
			"WhenVictoryConditionIsFinalBattle_ReturnsTrue",
			func(configuration *config.GeneratorConfig) {
				configuration.GameEndConditions.VictoryCondition = "win_condition_4"
			},
			true,
		},
		{
			"WhenGladiatorArenaRulesAreDisabled_ReturnsFalse",
			func(_ *config.GeneratorConfig) {},
			false,
		},
		{
			"WhenGladiatorArenaRulesAreNil_ReturnsFalse",
			func(configuration *config.GeneratorConfig) {
				configuration.GladiatorArenaRules = nil
			},
			false,
		},
		{
			"WhenGameEndConditionsAreNil_ReturnsFalse",
			func(configuration *config.GeneratorConfig) {
				configuration.GameEndConditions = nil
			},
			false,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.subtestName, func(t *testing.T) {
			t.Parallel()
			// Arrange
			configuration := config.NewGeneratorConfig()
			testCase.mutate(configuration)

			// Act
			actual := configuration.IsGladiatorArenaMode()

			// Assert
			assert.Equal(t, testCase.expected, actual)
		})
	}
}
