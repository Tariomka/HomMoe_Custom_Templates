package editorState_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template_entity"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenAnEffectiveModeChanges_ManualEditsAreCleared(t *testing.T) {
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
		{
			"WhenTheFinalBattleVictoryConditionIsChosen_ManualEditsAreCleared",
			func(state *editor_state_model.EditorState) {
				state.VictoryCondition = registry.GetWinningConditionValues().FinalBattle
			},
		},
		{
			"WhenBothModesAreEnabledAtOnce_ManualEditsAreCleared",
			func(state *editor_state_model.EditorState) {
				state.Tournament = true
				state.GladiatorArena = true
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.subtestName, func(t *testing.T) {
			t.Parallel()
			// Arrange
			previous := editor_state_model.NewDefaultEditorStateModel()
			candidate := previous.Clone()
			testCase.mutate(&candidate)
			candidate.PlayerCount = editor_state_model.TournamentPlayerCount
			withManualEdits(&candidate)
			requested := candidate.Clone()

			// Act
			candidate.ApplyModeTransition(&previous, &requested)

			// Assert
			assert.False(t, candidate.HasManualEdits())
		})
	}
}

func TestWhenAnEffectiveModeIsLeft_ManualEditsAreCleared(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		subtestName string
		mutate      func(state *editor_state_model.EditorState)
	}{
		{
			"WhenTournamentIsDisabled_ManualEditsAreCleared",
			func(state *editor_state_model.EditorState) { state.Tournament = false },
		},
		{
			"WhenTheArenaIsDisabled_ManualEditsAreCleared",
			func(state *editor_state_model.EditorState) { state.GladiatorArena = false },
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.subtestName, func(t *testing.T) {
			t.Parallel()
			// Arrange
			previous := editor_state_model.NewDefaultEditorStateModel()
			previous.Tournament = true
			previous.GladiatorArena = true
			previous.PlayerCount = editor_state_model.TournamentPlayerCount
			candidate := previous.Clone()
			testCase.mutate(&candidate)
			withManualEdits(&candidate)
			requested := candidate.Clone()

			// Act
			candidate.ApplyModeTransition(&previous, &requested)

			// Assert
			assert.False(t, candidate.HasManualEdits())
		})
	}
}

func TestWhenAnEffectiveModeChangesWithManualEdits_ReportsTheDiscard(t *testing.T) {
	t.Parallel()
	// Arrange
	previous := editor_state_model.NewDefaultEditorStateModel()
	candidate := previous.Clone()
	candidate.GladiatorArena = true
	withManualEdits(&candidate)
	requested := candidate.Clone()

	// Act
	outcome := candidate.ApplyModeTransition(&previous, &requested)

	// Assert
	assert.True(t, outcome.ManualEditsDiscarded)
}

func TestWhenAnEffectiveModeChangesWithoutManualEdits_ReportsNoDiscard(t *testing.T) {
	t.Parallel()
	// Arrange
	previous := editor_state_model.NewDefaultEditorStateModel()
	candidate := previous.Clone()
	candidate.GladiatorArena = true
	requested := candidate.Clone()

	// Act
	outcome := candidate.ApplyModeTransition(&previous, &requested)

	// Assert
	assert.False(t, outcome.ManualEditsDiscarded)
}

func TestWhenOnlyManualZonesExistOnAModeChange_TheyAreCleared(t *testing.T) {
	t.Parallel()
	// Arrange
	previous := editor_state_model.NewDefaultEditorStateModel()
	candidate := previous.Clone()
	candidate.Tournament = true
	candidate.PlayerCount = editor_state_model.TournamentPlayerCount
	candidate.ManualZones = []template_model.Zone{{Name: "Zone A"}}
	requested := candidate.Clone()

	// Act
	outcome := candidate.ApplyModeTransition(&previous, &requested)

	// Assert
	assert.True(t, outcome.ManualEditsDiscarded)
}

func TestWhenOnlyManualConnectionsExistOnAModeChange_TheyAreCleared(t *testing.T) {
	t.Parallel()
	// Arrange
	previous := editor_state_model.NewDefaultEditorStateModel()
	candidate := previous.Clone()
	candidate.Tournament = true
	candidate.PlayerCount = editor_state_model.TournamentPlayerCount
	candidate.ManualConnections = []editor_state_model.ManualConnectionSave{
		{Connection: template_entity.Connection{Name: "A-B"}},
	}
	requested := candidate.Clone()

	// Act
	outcome := candidate.ApplyModeTransition(&previous, &requested)

	// Assert
	assert.True(t, outcome.ManualEditsDiscarded)
}

// A tournament transition is judged on the mode, not on whether the count has
// to move: two players already is still a different generated map.
func TestWhenTournamentStartsAtTwoPlayers_ManualEditsAreStillCleared(t *testing.T) {
	t.Parallel()
	// Arrange
	previous := editor_state_model.NewDefaultEditorStateModel()
	previous.PlayerCount = editor_state_model.TournamentPlayerCount
	candidate := previous.Clone()
	candidate.Tournament = true
	withManualEdits(&candidate)
	requested := candidate.Clone()

	// Act
	outcome := candidate.ApplyModeTransition(&previous, &requested)

	// Assert
	assert.True(t, outcome.ManualEditsDiscarded)
}

func TestWhenModeRepresentationChangesButEffectiveModesDoNot_ManualEditsSurvive(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		subtestName string
		mutate      func(state *editor_state_model.EditorState)
	}{
		{
			"WhenTheCheckboxIsTradedForTheSelector_ManualEditsSurvive",
			func(state *editor_state_model.EditorState) {
				state.Tournament = false
				state.VictoryCondition = registry.GetWinningConditionValues().Tournament
			},
		},
		{
			"WhenTheSelectorIsTradedForTheCheckbox_ManualEditsSurvive",
			func(state *editor_state_model.EditorState) {
				state.Tournament = true
				state.VictoryCondition = registry.GetWinningConditionValues().Standard
			},
		},
		{
			"WhenOnlyOneOfTwoTournamentAliasesIsRemoved_ManualEditsSurvive",
			func(state *editor_state_model.EditorState) { state.Tournament = false },
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.subtestName, func(t *testing.T) {
			t.Parallel()
			// Arrange
			previous := editor_state_model.NewDefaultEditorStateModel()
			previous.Tournament = true
			previous.VictoryCondition = registry.GetWinningConditionValues().Tournament
			previous.PlayerCount = editor_state_model.TournamentPlayerCount
			candidate := previous.Clone()
			testCase.mutate(&candidate)
			withManualEdits(&candidate)
			requested := candidate.Clone()

			// Act
			outcome := candidate.ApplyModeTransition(&previous, &requested)

			// Assert
			assert.False(t, outcome.ManualEditsDiscarded)
		})
	}
}

func TestWhenOnlyUnrelatedRulesChange_ManualEditsSurvive(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		subtestName string
		mutate      func(state *editor_state_model.EditorState)
	}{
		{
			"WhenTheVictoryConditionChangesToAnUnrelatedOne_ManualEditsSurvive",
			func(state *editor_state_model.EditorState) {
				state.VictoryCondition = registry.GetWinningConditionValues().CapitalCapture
			},
		},
		{
			"WhenArenaTimingChanges_ManualEditsSurvive",
			func(state *editor_state_model.EditorState) { state.GladiatorArenaCountDay++ },
		},
		{
			"WhenTournamentTimingChanges_ManualEditsSurvive",
			func(state *editor_state_model.EditorState) { state.TournamentInterval++ },
		},
		{
			"WhenTournamentPointsChange_ManualEditsSurvive",
			func(state *editor_state_model.EditorState) { state.TournamentPointsToWin++ },
		},
		{
			"WhenSaveArmyChanges_ManualEditsSurvive",
			func(state *editor_state_model.EditorState) { state.TournamentSaveArmy = !state.TournamentSaveArmy },
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.subtestName, func(t *testing.T) {
			t.Parallel()
			// Arrange
			previous := editor_state_model.NewDefaultEditorStateModel()
			candidate := previous.Clone()
			testCase.mutate(&candidate)
			withManualEdits(&candidate)
			requested := candidate.Clone()

			// Act
			candidate.ApplyModeTransition(&previous, &requested)

			// Assert
			assert.True(t, candidate.HasManualEdits())
		})
	}
}

func TestWhenALoadCorrectsTheTournamentPlayerCount_ReportsTheCorrection(t *testing.T) {
	t.Parallel()
	// Arrange
	requested := editor_state_model.NewDefaultEditorStateModel()
	requested.Tournament = true
	requested.PlayerCount = gofakeit.Number(3, 8)
	candidate := requested.Clone()
	candidate.PlayerCount = editor_state_model.TournamentPlayerCount

	// Act
	outcome := candidate.ApplyModeTransition(nil, &requested)

	// Assert
	assert.True(t, outcome.TournamentCountCorrected)
}

func TestWhenALoadCorrectsTheTournamentPlayerCount_ManualEditsAreCleared(t *testing.T) {
	t.Parallel()
	// Arrange
	requested := editor_state_model.NewDefaultEditorStateModel()
	requested.VictoryCondition = registry.GetWinningConditionValues().Tournament
	requested.PlayerCount = gofakeit.Number(3, 8)
	withManualEdits(&requested)
	candidate := requested.Clone()
	candidate.PlayerCount = editor_state_model.TournamentPlayerCount

	// Act
	outcome := candidate.ApplyModeTransition(nil, &requested)

	// Assert
	assert.True(t, outcome.ManualEditsDiscarded)
}

func TestWhenALoadedTournamentAlreadyHasTwoPlayers_ManualEditsSurvive(t *testing.T) {
	t.Parallel()
	// Arrange
	requested := editor_state_model.NewDefaultEditorStateModel()
	requested.Tournament = true
	requested.PlayerCount = editor_state_model.TournamentPlayerCount
	withManualEdits(&requested)
	candidate := requested.Clone()

	// Act
	outcome := candidate.ApplyModeTransition(nil, &requested)

	// Assert
	assert.False(t, outcome.ManualEditsDiscarded)
}

func TestWhenALoadedTournamentAlreadyHasTwoPlayers_ReportsNoCorrection(t *testing.T) {
	t.Parallel()
	// Arrange
	requested := editor_state_model.NewDefaultEditorStateModel()
	requested.Tournament = true
	requested.PlayerCount = editor_state_model.TournamentPlayerCount
	candidate := requested.Clone()

	// Act
	outcome := candidate.ApplyModeTransition(nil, &requested)

	// Assert
	assert.False(t, outcome.TournamentCountCorrected)
}

// A count clamped outside tournament mode is an ordinary range fix, not the
// tournament correction the user has to be told about.
func TestWhenANonTournamentCountIsClamped_ReportsNoCorrection(t *testing.T) {
	t.Parallel()
	// Arrange
	requested := editor_state_model.NewDefaultEditorStateModel()
	requested.PlayerCount = gofakeit.Number(9, 100)
	candidate := requested.Clone()
	candidate.PlayerCount = 8

	// Act
	outcome := candidate.ApplyModeTransition(nil, &requested)

	// Assert
	assert.False(t, outcome.TournamentCountCorrected)
}

func TestWhenALoadKeepsEveryModeValid_ManualEditsSurvive(t *testing.T) {
	t.Parallel()
	// Arrange
	requested := editor_state_model.NewDefaultEditorStateModel()
	requested.GladiatorArena = true
	withManualEdits(&requested)
	candidate := requested.Clone()

	// Act
	candidate.ApplyModeTransition(nil, &requested)

	// Assert
	assert.True(t, candidate.HasManualEdits())
}

func TestWhenATransitionIsApplied_ThePreviousStateIsUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	previous := editor_state_model.NewDefaultEditorStateModel()
	withManualEdits(&previous)
	untouched := previous.Clone()
	candidate := previous.Clone()
	candidate.Tournament = true
	candidate.PlayerCount = editor_state_model.TournamentPlayerCount
	requested := candidate.Clone()

	// Act
	candidate.ApplyModeTransition(&previous, &requested)

	// Assert
	assert.Equal(t, untouched, previous)
}

func TestWhenATransitionIsApplied_TheRequestedStateIsUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	previous := editor_state_model.NewDefaultEditorStateModel()
	candidate := previous.Clone()
	candidate.Tournament = true
	candidate.PlayerCount = editor_state_model.TournamentPlayerCount
	withManualEdits(&candidate)
	requested := candidate.Clone()
	untouched := requested.Clone()

	// Act
	candidate.ApplyModeTransition(&previous, &requested)

	// Assert
	assert.Equal(t, untouched, requested)
}

// withManualEdits stamps a zone and a connection onto a state, standing in for
// a hand-made layout.
func withManualEdits(state *editor_state_model.EditorState) {
	state.ManualZones = []template_model.Zone{{Name: "Zone A"}}
	state.ManualConnections = []editor_state_model.ManualConnectionSave{
		{Connection: template_entity.Connection{Name: "A-B"}},
	}
}
