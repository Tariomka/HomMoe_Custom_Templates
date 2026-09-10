package integration_test

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"os"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/composition"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenTournamentSaveArmyIsConfigured_SerializesExpectedFieldPresence(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name                  string
		activateTournament    func(state *editor_state_model.EditorState)
		expectedSaveArmyValue jsontext.Value
	}{
		{
			name: "WhenTournamentCheckboxEnabledAndSaveArmyEnabled_SaveArmyIsSerialized",
			activateTournament: func(state *editor_state_model.EditorState) {
				state.Tournament = true
			},
			expectedSaveArmyValue: jsontext.Value("true"),
		},
		{
			name: "WhenTournamentCheckboxEnabledAndSaveArmyDisabled_SaveArmyIsOmitted",
			activateTournament: func(state *editor_state_model.EditorState) {
				state.Tournament = true
			},
		},
		{
			name: "WhenTournamentVictorySelectorSelectedAndSaveArmyEnabled_SaveArmyIsSerialized",
			activateTournament: func(state *editor_state_model.EditorState) {
				state.VictoryCondition = registry.GetWinningConditionValues().Tournament
			},
			expectedSaveArmyValue: jsontext.Value("true"),
		},
		{
			name: "WhenTournamentVictorySelectorSelectedAndSaveArmyDisabled_SaveArmyIsOmitted",
			activateTournament: func(state *editor_state_model.EditorState) {
				state.VictoryCondition = registry.GetWinningConditionValues().Tournament
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			// Arrange
			handler := composition.InitializeGuiHandler()
			state := editor_state_model.NewDefaultEditorStateModel()
			state.TemplateName = "Tournament Save Army"
			state.TournamentSaveArmy = testCase.expectedSaveArmyValue != nil
			testCase.activateTournament(&state)

			// Act
			generated, err := handler.GenerateTemplate(editor_state_dto.EditorStateDto{EditorState: state})
			require.NoError(t, err)
			savedPath, err := handler.SaveTemplate(dtos.TemplateSaveDto{
				Template:   generated.Template,
				Topology:   state.Topology,
				OutputPath: t.TempDir(),
			})
			require.NoError(t, err)
			rawTemplate, err := os.ReadFile(savedPath)

			// Assert
			require.NoError(t, err)
			var template map[string]jsontext.Value
			require.NoError(t, json.Unmarshal(rawTemplate, &template))
			gameRulesValue, found := template["gameRules"]
			require.True(t, found)
			var gameRules map[string]jsontext.Value
			require.NoError(t, json.Unmarshal(gameRulesValue, &gameRules))
			winConditionsValue, found := gameRules["winConditions"]
			require.True(t, found)
			var winConditions map[string]jsontext.Value
			require.NoError(t, json.Unmarshal(winConditionsValue, &winConditions))
			var tournament bool
			require.NoError(t, json.Unmarshal(winConditions["tournament"], &tournament))
			require.True(t, tournament)
			assert.Equal(t, testCase.expectedSaveArmyValue, winConditions["tournamentSaveArmy"])
		})
	}
}
