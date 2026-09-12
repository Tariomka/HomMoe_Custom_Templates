package state_test

import (
	"errors"
	"testing"
	"time"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/drivers"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestWhenUpdateChangesState_ChangeIsApplied(t *testing.T) {
	t.Parallel()
	// Arrange
	state := drivers.NewUIState(
		&test_helpers.TemplateHandlerMock{},
		test_helpers.NewFileSystemHandler(),
		test_helpers.NewRegenerationHandler(),

		false)

	playerCount := gofakeit.Number(3, 8)

	// Act
	state.UpdateState(func(dto *editor_state_model.EditorState) { dto.PlayerCount = playerCount })

	// Assert
	assert.Equal(t, playerCount, state.GetStateData().PlayerCount)
}

func TestWhenUpdateChangesStateAfterGeneration_StateBecomesUnsaved(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newGeneratedState()

	// Act
	state.UpdateState(func(dto *editor_state_model.EditorState) { dto.TemplateName = gofakeit.ProductName() })

	// Assert
	assert.True(t, state.IsUnsaved())
}

func TestWhenUpdateDoesNotChangeStateAfterGeneration_StateStaysSaved(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newGeneratedState()

	// Act
	state.UpdateState(func(_ *editor_state_model.EditorState) {})

	// Assert
	assert.False(t, state.IsUnsaved())
}

// A discarded manual layout is a persisted change, and the change detection
// deliberately ignores the manual fields, so the discard has to flag the
// document by itself - here there is not even a generation to compare against.
func TestWhenAModeChangeDiscardsManualEdits_StateBecomesUnsaved(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newStateWithManualEdits()
	require.False(t, state.IsUnsaved(), "precondition: manual edits alone do not flag an ungenerated document")

	// Act
	state.UpdateState(func(dto *editor_state_model.EditorState) { dto.GladiatorArena = true })

	// Assert
	assert.True(t, state.IsUnsaved())
}

func TestWhenAModeChangeDiscardsManualEdits_ExitWarnsAgain(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newStateWithManualEdits()
	exitCalled := false
	state.SetOnExit(func() { exitCalled = true })

	// Act
	state.UpdateState(func(dto *editor_state_model.EditorState) { dto.GladiatorArena = true })
	state.Exit()

	// Assert
	assert.False(t, exitCalled)
}

func TestWhenAModeChangeDiscardsManualEdits_ManualEditsAreGone(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newStateWithManualEdits()

	// Act
	state.UpdateState(func(dto *editor_state_model.EditorState) { dto.Tournament = true })

	// Assert
	stateData := state.GetStateData()
	assert.False(t, stateData.HasManualEdits())
}

func TestWhenAModeChangeDiscardsManualEdits_StatusReportsTheDiscard(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newStateWithManualEdits()

	// Act
	state.UpdateState(func(dto *editor_state_model.EditorState) { dto.Tournament = true })

	// Assert
	message, _ := state.GetStatus()
	assert.Contains(t, message, "manual zone layout was discarded")
}

// The mode change regenerates on the same frame, and that generation writes its
// own status; the warning has to survive it.
func TestWhenGenerationFollowsADiscard_TheGeneratedStatusCarriesTheNotice(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newStateWithManualEdits()
	state.UpdateState(func(dto *editor_state_model.EditorState) { dto.Tournament = true })

	// Act
	state.Generate()

	// Assert
	message, _ := state.GetStatus()
	assert.Contains(t, message, "manual zone layout was discarded")
}

func TestWhenASecondGenerationFollowsADiscard_TheNoticeIsNotRepeated(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newStateWithManualEdits()
	state.UpdateState(func(dto *editor_state_model.EditorState) { dto.Tournament = true })
	state.Generate()

	// Act
	state.Generate()

	// Assert
	message, _ := state.GetStatus()
	assert.NotContains(t, message, "manual zone layout was discarded")
}

// A failed generation must not swallow the warning: the retry still owes it.
func TestWhenGenerationFailsAfterADiscard_TheRetryCarriesTheNotice(t *testing.T) {
	t.Parallel()
	// Arrange
	handlerMock := &test_helpers.TemplateHandlerMock{}
	template := test_helpers.GetDefaultTemplateModel()
	handlerMock.On("GenerateTemplate", mock.Anything).
		Return(dtos.TemplateLoadDto{}, errors.New("generator unavailable")).Once()
	handlerMock.On("GenerateTemplate", mock.Anything).Return(dtos.TemplateLoadDto{Template: &template}, nil)
	state := newStateWithHandler(handlerMock)
	state.UpdateState(func(dto *editor_state_model.EditorState) {
		dto.ManualZones = []template_model.Zone{{Name: "Zone A"}}
	})
	state.UpdateState(func(dto *editor_state_model.EditorState) { dto.Tournament = true })
	state.Generate()

	// Act
	state.Generate()

	// Assert
	message, _ := state.GetStatus()
	assert.Contains(t, message, "manual zone layout was discarded")
}

func TestWhenTheDocumentIsReset_APendingNoticeIsDropped(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newStateWithManualEdits()
	state.UpdateState(func(dto *editor_state_model.EditorState) { dto.Tournament = true })

	// Act
	state.Reset()
	state.Generate()

	// Assert
	message, _ := state.GetStatus()
	assert.NotContains(t, message, "manual zone layout was discarded")
}

func TestWhenIdleUpdatesFollowFailedGeneration_TheErrorStatusSurvives(t *testing.T) {
	t.Parallel()
	// Arrange
	handlerMock := &test_helpers.TemplateHandlerMock{}
	handlerMock.On("GenerateTemplate", mock.Anything).
		Return(dtos.TemplateLoadDto{}, errors.New("generator unavailable"))
	state := newStateWithHandler(handlerMock)
	state.UpdateState(func(candidate *editor_state_model.EditorState) {
		candidate.ManualZones = []template_model.Zone{{Name: "Zone A"}}
	})
	state.UpdateState(func(candidate *editor_state_model.EditorState) { candidate.Tournament = true })
	state.AutoRegenerate(time.Now())

	// Act
	for range 5 {
		state.UpdateState(func(_ *editor_state_model.EditorState) {})
		state.AutoRegenerate(time.Now())
	}

	// Assert
	message, isError := state.GetStatus()
	assert.Equal(t, []any{"Generation failed: generator unavailable.", true}, []any{message, isError})
}

// A corrected player count is a change the user did not ask for and has not
// saved, and before the first generation there is no snapshot for the change
// detection to notice it with.
func TestWhenAnUpdateCorrectsTheTournamentCount_StateBecomesUnsaved(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newStateWithHandler(newTournamentCorrectingHandler())
	require.False(t, state.IsUnsaved(), "precondition: a fresh document is saved")

	// Act
	state.UpdateState(func(dto *editor_state_model.EditorState) {
		dto.Tournament = true
		dto.PlayerCount = gofakeit.Number(3, 8)
	})

	// Assert
	assert.True(t, state.IsUnsaved())
}

func TestWhenAnUpdateCorrectsTheTournamentCount_ExitWarnsAgain(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newStateWithHandler(newTournamentCorrectingHandler())
	exitCalled := false
	state.SetOnExit(func() { exitCalled = true })

	// Act
	state.UpdateState(func(dto *editor_state_model.EditorState) {
		dto.Tournament = true
		dto.PlayerCount = gofakeit.Number(3, 8)
	})
	state.Exit()

	// Assert
	assert.False(t, exitCalled)
}

func TestWhenAnUpdateCorrectsTheTournamentCount_TheStatusReportsIt(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newStateWithHandler(newTournamentCorrectingHandler())

	// Act
	state.UpdateState(func(dto *editor_state_model.EditorState) {
		dto.Tournament = true
		dto.PlayerCount = gofakeit.Number(3, 8)
	})

	// Assert
	message, _ := state.GetStatus()
	assert.Contains(t, message, "player count was corrected")
}

// Two transitions can land before the generation that reports them; the second
// one must not silence the first.
func TestWhenASecondTransitionPrecedesGeneration_TheDiscardIsStillReported(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newStateWithHandler(newTournamentCorrectingHandler())
	state.UpdateState(func(dto *editor_state_model.EditorState) {
		dto.ManualZones = []template_model.Zone{{Name: "Zone A"}}
	})
	state.UpdateState(func(dto *editor_state_model.EditorState) { dto.GladiatorArena = true })

	// Act
	state.UpdateState(func(dto *editor_state_model.EditorState) {
		dto.Tournament = true
		dto.PlayerCount = gofakeit.Number(3, 8)
	})
	state.Generate()

	// Assert
	message, _ := state.GetStatus()
	assert.Contains(t, message, "manual zone layout was discarded")
}

func TestWhenASecondTransitionPrecedesGeneration_TheCorrectionIsAlsoReported(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newStateWithHandler(newTournamentCorrectingHandler())
	state.UpdateState(func(dto *editor_state_model.EditorState) {
		dto.ManualZones = []template_model.Zone{{Name: "Zone A"}}
	})
	state.UpdateState(func(dto *editor_state_model.EditorState) { dto.GladiatorArena = true })

	// Act
	state.UpdateState(func(dto *editor_state_model.EditorState) {
		dto.Tournament = true
		dto.PlayerCount = gofakeit.Number(3, 8)
	})
	state.Generate()

	// Assert
	message, _ := state.GetStatus()
	assert.Contains(t, message, "player count was corrected")
}

// newTournamentCorrectingHandler returns a handler that generates, and that
// holds a tournament state at two players the way the real validator does.
func newTournamentCorrectingHandler() *test_helpers.TemplateHandlerMock {
	handlerMock := &test_helpers.TemplateHandlerMock{}
	template := test_helpers.GetDefaultTemplateModel()
	handlerMock.On("GenerateTemplate", mock.Anything).Return(dtos.TemplateLoadDto{Template: &template}, nil)
	handlerMock.ValidateEditorStateFunc = func(
		state editor_state_model.EditorState,
		_ bool,
	) editor_state_dto.EditorStateValidationDto {
		if state.IsEffectiveTournament() {
			state.PlayerCount = editor_state_model.TournamentPlayerCount
		}
		return editor_state_dto.EditorStateValidationDto{State: state}
	}
	return handlerMock
}

// newGeneratedState returns a State that has generated once, so a previous
// state snapshot exists and change detection is active.
func newGeneratedState() *drivers.State {
	handlerMock := &test_helpers.TemplateHandlerMock{}
	template := test_helpers.GetDefaultTemplateModel()
	handlerMock.On("GenerateTemplate", mock.Anything).Return(dtos.TemplateLoadDto{Template: &template}, nil)
	state := newStateWithHandler(handlerMock)

	state.Generate()
	return state
}

// newStateWithManualEdits returns an ungenerated State carrying a hand-made
// layout, which is what a mode change has to discard.
func newStateWithManualEdits() *drivers.State {
	handlerMock := &test_helpers.TemplateHandlerMock{}
	template := test_helpers.GetDefaultTemplateModel()
	handlerMock.On("GenerateTemplate", mock.Anything).Return(dtos.TemplateLoadDto{Template: &template}, nil)
	state := newStateWithHandler(handlerMock)
	state.UpdateState(func(dto *editor_state_model.EditorState) {
		dto.ManualZones = []template_model.Zone{{Name: "Zone A"}}
	})
	return state
}

func newStateWithHandler(handlerMock *test_helpers.TemplateHandlerMock) *drivers.State {
	return drivers.NewUIState(
		handlerMock,
		test_helpers.NewFileSystemHandler(),
		test_helpers.NewRegenerationHandler(),

		false)
}
