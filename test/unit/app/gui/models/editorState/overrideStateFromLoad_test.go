package editorState_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenALoadedTournamentCountIsCorrected_ReportsTheCorrection(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newLoadCorrectingEditorState()
	raw := editor_state_model.NewDefaultEditorStateModel()
	raw.Tournament = true
	raw.PlayerCount = gofakeit.Number(3, 8)
	corrected := raw.Clone()
	corrected.PlayerCount = editor_state_model.TournamentPlayerCount

	// Act
	outcome := state.OverrideStateFromLoad(raw, corrected)

	// Assert
	assert.True(t, outcome.TournamentCountCorrected)
}

func TestWhenALoadedTournamentCountIsCorrected_ManualEditsAreDropped(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newLoadCorrectingEditorState()
	raw := editor_state_model.NewDefaultEditorStateModel()
	raw.VictoryCondition = registry.GetWinningConditionValues().Tournament
	raw.PlayerCount = 6
	raw.ManualZones = []template_model.Zone{{Name: "Zone A"}}
	corrected := raw.Clone()
	corrected.PlayerCount = editor_state_model.TournamentPlayerCount

	// Act
	state.OverrideStateFromLoad(raw, corrected)

	// Assert
	assert.False(t, state.HasManualEdits())
}

func TestWhenALoadedTournamentCountIsCorrected_ReportsTheDiscard(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newLoadCorrectingEditorState()
	raw := editor_state_model.NewDefaultEditorStateModel()
	raw.Tournament = true
	raw.PlayerCount = 6
	raw.ManualZones = []template_model.Zone{{Name: "Zone A"}}
	corrected := raw.Clone()
	corrected.PlayerCount = editor_state_model.TournamentPlayerCount

	// Act
	outcome := state.OverrideStateFromLoad(raw, corrected)

	// Assert
	assert.True(t, outcome.ManualEditsDiscarded)
}

func TestWhenAValidTournamentIsLoaded_ManualEditsSurvive(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newLoadCorrectingEditorState()
	raw := editor_state_model.NewDefaultEditorStateModel()
	raw.Tournament = true
	raw.PlayerCount = editor_state_model.TournamentPlayerCount
	raw.ManualZones = []template_model.Zone{{Name: "Zone A"}}
	corrected := raw.Clone()

	// Act
	state.OverrideStateFromLoad(raw, corrected)

	// Assert
	assert.True(t, state.HasManualEdits())
}

// A load replaces the document, so the modes of the document being replaced are
// never compared against - only the file's own correction matters.
func TestWhenALoadChangesTheEffectiveMode_ManualEditsSurvive(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newLoadCorrectingEditorState()
	state.UpdateCurrentState(func(current *editor_state_model.EditorState) { current.GladiatorArena = true })
	raw := editor_state_model.NewDefaultEditorStateModel()
	raw.ManualZones = []template_model.Zone{{Name: "Zone A"}}
	corrected := raw.Clone()

	// Act
	outcome := state.OverrideStateFromLoad(raw, corrected)

	// Assert
	assert.False(t, outcome.ManualEditsDiscarded)
}

func TestWhenAStateIsLoaded_TheStoredStateDoesNotAliasTheCaller(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newLoadCorrectingEditorState()
	raw := editor_state_model.NewDefaultEditorStateModel()
	raw.ManualZones = []template_model.Zone{{Name: "Zone A"}}
	corrected := raw.Clone()
	state.OverrideStateFromLoad(raw, corrected)

	// Act
	corrected.ManualZones[0].Name = "Renamed"

	// Assert
	assert.Equal(t, "Zone A", state.GetManualZones()[0].Name)
}

// newLoadCorrectingEditorState returns a state whose validation corrects a
// tournament player count, as the real validator does.
func newLoadCorrectingEditorState() *models.EditorState {
	return newEditorStateWithValidation(
		func(stateDto editor_state_model.EditorState, _ bool) editor_state_dto.EditorStateValidationDto {
			if stateDto.IsEffectiveTournament() {
				stateDto.PlayerCount = editor_state_model.TournamentPlayerCount
			}
			return editor_state_dto.EditorStateValidationDto{State: stateDto}
		},
	)
}
