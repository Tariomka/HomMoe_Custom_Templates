package editorState_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
	"github.com/stretchr/testify/assert"
)

func TestWhenArenaCheckboxIsSet_ReportsEffectiveGladiatorArena(t *testing.T) {
	t.Parallel()
	// Arrange
	state := editor_state_model.NewDefaultEditorStateModel()
	state.GladiatorArena = true

	// Act
	effective := state.IsEffectiveGladiatorArena()

	// Assert
	assert.True(t, effective)
}

// Guardian Arena is the UI label of the final battle victory condition.
func TestWhenVictoryConditionIsFinalBattle_ReportsEffectiveGladiatorArena(t *testing.T) {
	t.Parallel()
	// Arrange
	state := editor_state_model.NewDefaultEditorStateModel()
	state.VictoryCondition = registry.GetWinningConditionValues().FinalBattle

	// Act
	effective := state.IsEffectiveGladiatorArena()

	// Assert
	assert.True(t, effective)
}

func TestWhenNeitherArenaAliasIsSet_ReportsNoEffectiveGladiatorArena(t *testing.T) {
	t.Parallel()
	// Arrange
	state := editor_state_model.NewDefaultEditorStateModel()

	// Act
	effective := state.IsEffectiveGladiatorArena()

	// Assert
	assert.False(t, effective)
}

func TestWhenArenaAliasMatchesConfig_PredicatesAgree(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		subtestName  string
		mutateState  func(state *editor_state_model.EditorState)
		mutateConfig func(generatorConfig *config.GeneratorConfig)
	}{
		{
			"WhenNothingEnablesTheArena_BothReportFalse",
			func(_ *editor_state_model.EditorState) {},
			func(_ *config.GeneratorConfig) {},
		},
		{
			"WhenTheRuleIsEnabled_BothReportTrue",
			func(state *editor_state_model.EditorState) { state.GladiatorArena = true },
			func(generatorConfig *config.GeneratorConfig) { generatorConfig.GladiatorArenaRules.Enabled = true },
		},
		{
			"WhenTheVictoryConditionIsFinalBattle_BothReportTrue",
			func(state *editor_state_model.EditorState) {
				state.VictoryCondition = registry.GetWinningConditionValues().FinalBattle
			},
			func(generatorConfig *config.GeneratorConfig) {
				generatorConfig.GameEndConditions.VictoryCondition = registry.GetWinningConditionValues().FinalBattle
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.subtestName, func(t *testing.T) {
			t.Parallel()
			// Arrange
			state := editor_state_model.NewDefaultEditorStateModel()
			generatorConfig := config.NewGeneratorConfig()
			testCase.mutateState(&state)
			testCase.mutateConfig(generatorConfig)

			// Act
			effective := state.IsEffectiveGladiatorArena()

			// Assert
			assert.Equal(t, generatorConfig.IsGladiatorArenaMode(), effective)
		})
	}
}
