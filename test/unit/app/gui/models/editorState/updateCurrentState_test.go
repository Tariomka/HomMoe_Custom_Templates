package editorState_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenUpdateChangesPlayerCount_ChangeIsApplied(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()
	playerCount := gofakeit.Number(3, 8)

	// Act
	state.UpdateCurrentState(func(dto *editor_state_model.EditorState) { dto.PlayerCount = playerCount })

	// Assert
	assert.Equal(t, playerCount, state.GetCurrentState().PlayerCount)
}

func TestWhenUpdateSetsPlayerCountAboveMaximum_PlayerCountIsClamped(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorStateWithValidation(func(
		stateDto editor_state_model.EditorState,
		_ bool,
	) editor_state_dto.EditorStateValidationDto {
		stateDto.PlayerCount = 8
		return editor_state_dto.EditorStateValidationDto{State: stateDto}
	})
	tooManyPlayers := gofakeit.Number(9, 100)

	// Act
	state.UpdateCurrentState(func(dto *editor_state_model.EditorState) { dto.PlayerCount = tooManyPlayers })

	// Assert
	assert.Equal(t, 8, state.GetCurrentState().PlayerCount)
}

func TestWhenUpdateSetsUnknownGameMode_GameModeIsResetToClassic(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorStateWithValidation(func(
		stateDto editor_state_model.EditorState,
		_ bool,
	) editor_state_dto.EditorStateValidationDto {
		stateDto.GameMode = registry.GetGameModeValues().Classic
		return editor_state_dto.EditorStateValidationDto{State: stateDto}
	})

	// Act
	state.UpdateCurrentState(func(dto *editor_state_model.EditorState) { dto.GameMode = "NotARealGameMode" })

	// Assert
	assert.Equal(t, registry.GetGameModeValues().Classic, state.GetCurrentState().GameMode)
}

func TestWhenUpdateChangesAnEffectiveMode_ManualEditsAreCleared(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		subtestName string
		mutate      func(state *editor_state_model.EditorState)
	}{
		{
			"WhenTournamentIsEnabled_ManualEditsAreCleared",
			func(state *editor_state_model.EditorState) { state.Tournament = true },
		},
		{
			"WhenTheTournamentVictoryConditionIsChosen_ManualEditsAreCleared",
			func(state *editor_state_model.EditorState) {
				state.VictoryCondition = registry.GetWinningConditionValues().Tournament
			},
		},
		{
			"WhenTheArenaIsEnabled_ManualEditsAreCleared",
			func(state *editor_state_model.EditorState) { state.GladiatorArena = true },
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.subtestName, func(t *testing.T) {
			t.Parallel()
			// Arrange
			state := newEditorState()
			state.UpdateCurrentState(func(current *editor_state_model.EditorState) {
				current.ManualZones = []template_model.Zone{{Name: "Zone A"}}
			})

			// Act
			state.UpdateCurrentState(testCase.mutate)

			// Assert
			assert.False(t, state.HasManualEdits())
		})
	}
}

func TestWhenUpdateChangesAnEffectiveModeWithManualEdits_ReportsTheDiscard(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()
	state.UpdateCurrentState(func(current *editor_state_model.EditorState) {
		current.ManualZones = []template_model.Zone{{Name: "Zone A"}}
	})

	// Act
	outcome := state.UpdateCurrentState(func(current *editor_state_model.EditorState) {
		current.GladiatorArena = true
	})

	// Assert
	assert.True(t, outcome.ManualEditsDiscarded)
}

func TestWhenUpdateLeavesTheEffectiveModes_ManualEditsSurvive(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()
	state.UpdateCurrentState(func(current *editor_state_model.EditorState) {
		current.Tournament = true
		current.ManualZones = []template_model.Zone{{Name: "Zone A"}}
	})

	// Act
	outcome := state.UpdateCurrentState(func(current *editor_state_model.EditorState) {
		current.VictoryCondition = registry.GetWinningConditionValues().Tournament
		current.Tournament = false
	})

	// Assert
	assert.False(t, outcome.ManualEditsDiscarded)
}

func TestWhenUpdateChangesNothingRelevant_ReportsNoOutcome(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()

	// Act
	outcome := state.UpdateCurrentState(func(current *editor_state_model.EditorState) {
		current.TemplateName = gofakeit.ProductName()
	})

	// Assert
	assert.Equal(t, editor_state_model.ModeTransitionOutcome{}, outcome)
}

// The correction has to be judged against what the caller asked for, which the
// validated state no longer shows.
func TestWhenUpdateRequestsAWrongTournamentCount_ReportsTheCorrection(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorStateWithValidation(func(
		stateDto editor_state_model.EditorState,
		_ bool,
	) editor_state_dto.EditorStateValidationDto {
		if stateDto.IsEffectiveTournament() {
			stateDto.PlayerCount = editor_state_model.TournamentPlayerCount
		}
		return editor_state_dto.EditorStateValidationDto{State: stateDto}
	})

	// Act
	outcome := state.UpdateCurrentState(func(current *editor_state_model.EditorState) {
		current.Tournament = true
		current.PlayerCount = gofakeit.Number(3, 8)
	})

	// Assert
	assert.True(t, outcome.TournamentCountCorrected)
}
