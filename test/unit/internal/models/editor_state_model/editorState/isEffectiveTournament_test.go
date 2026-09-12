package editorState_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
	"github.com/stretchr/testify/assert"
)

func TestWhenTournamentCheckboxIsSet_ReportsEffectiveTournament(t *testing.T) {
	t.Parallel()
	// Arrange
	state := editor_state_model.NewDefaultEditorStateModel()
	state.Tournament = true

	// Act
	effective := state.IsEffectiveTournament()

	// Assert
	assert.True(t, effective)
}

func TestWhenVictoryConditionIsTournament_ReportsEffectiveTournament(t *testing.T) {
	t.Parallel()
	// Arrange
	state := editor_state_model.NewDefaultEditorStateModel()
	state.VictoryCondition = registry.GetWinningConditionValues().Tournament

	// Act
	effective := state.IsEffectiveTournament()

	// Assert
	assert.True(t, effective)
}

func TestWhenNeitherTournamentAliasIsSet_ReportsNoEffectiveTournament(t *testing.T) {
	t.Parallel()
	// Arrange
	state := editor_state_model.NewDefaultEditorStateModel()

	// Act
	effective := state.IsEffectiveTournament()

	// Assert
	assert.False(t, effective)
}

// The generator judges the very same rules from its own configuration, so both
// aliases must mean to the editor state what they mean to that configuration.
func TestWhenTournamentAliasMatchesConfig_PredicatesAgree(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		subtestName  string
		mutateState  func(state *editor_state_model.EditorState)
		mutateConfig func(generatorConfig *config.GeneratorConfig)
	}{
		{
			"WhenNothingEnablesTournament_BothReportFalse",
			func(_ *editor_state_model.EditorState) {},
			func(_ *config.GeneratorConfig) {},
		},
		{
			"WhenTheRuleIsEnabled_BothReportTrue",
			func(state *editor_state_model.EditorState) { state.Tournament = true },
			func(generatorConfig *config.GeneratorConfig) { generatorConfig.TournamentRules.Enabled = true },
		},
		{
			"WhenTheVictoryConditionIsTournament_BothReportTrue",
			func(state *editor_state_model.EditorState) {
				state.VictoryCondition = registry.GetWinningConditionValues().Tournament
			},
			func(generatorConfig *config.GeneratorConfig) {
				generatorConfig.GameEndConditions.VictoryCondition = registry.GetWinningConditionValues().Tournament
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
			effective := state.IsEffectiveTournament()

			// Assert
			assert.Equal(t, generatorConfig.IsTournamentMode(), effective)
		})
	}
}
