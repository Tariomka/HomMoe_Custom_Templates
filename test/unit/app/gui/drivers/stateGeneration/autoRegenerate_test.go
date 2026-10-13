package stateGeneration_test

import (
	"testing"
	"time"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/drivers"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers/handler_interfaces"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWhenStateWasNeverGenerated_GeneratesImmediately(t *testing.T) {
	t.Parallel()
	// Arrange
	state, handlerMock := newAutoRegenerateState()

	// Act
	state.AutoRegenerate(time.Now())

	// Assert
	handlerMock.AssertNumberOfCalls(t, "GenerateTemplate", 1)
}

func TestWhenStateIsUnchangedSinceGeneration_DoesNotRegenerate(t *testing.T) {
	t.Parallel()
	// Arrange
	state, handlerMock := newAutoRegenerateState()
	state.AutoRegenerate(time.Now())

	// Act
	state.AutoRegenerate(time.Now())

	// Assert
	handlerMock.AssertNumberOfCalls(t, "GenerateTemplate", 1)
}

func TestWhenStateIsUnchangedSinceGeneration_NoRedrawIsScheduled(t *testing.T) {
	t.Parallel()
	// Arrange
	state, _ := newAutoRegenerateState()
	state.AutoRegenerate(time.Now())

	// Act
	_, scheduleRedraw := state.AutoRegenerate(time.Now())

	// Assert
	assert.False(t, scheduleRedraw)
}

func TestWhenLayoutOptionChanges_RegeneratesImmediately(t *testing.T) {
	t.Parallel()
	// Arrange
	state, handlerMock := newAutoRegenerateState()
	state.AutoRegenerate(time.Now())
	state.UpdateState(func(dto *editor_state_model.EditorState) { dto.PlayerCount++ })

	// Act
	state.AutoRegenerate(time.Now())

	// Assert
	handlerMock.AssertNumberOfCalls(t, "GenerateTemplate", 2)
}

func TestWhenNonLayoutOptionChanges_DebounceTimerIsArmed(t *testing.T) {
	t.Parallel()
	// Arrange
	state, _ := newAutoRegenerateState()
	now := time.Now()
	state.AutoRegenerate(now)
	state.UpdateState(func(dto *editor_state_model.EditorState) { dto.TemplateName = gofakeit.ProductName() })

	// Act
	redrawAt, scheduleRedraw := state.AutoRegenerate(now)

	// Assert
	assert.Equal(t,
		[]any{now.Add(300 * time.Millisecond), true},
		[]any{redrawAt, scheduleRedraw})
}

func TestWhenNonLayoutOptionChanges_NoImmediateRegeneration(t *testing.T) {
	t.Parallel()
	// Arrange
	state, handlerMock := newAutoRegenerateState()
	now := time.Now()
	state.AutoRegenerate(now)
	state.UpdateState(func(dto *editor_state_model.EditorState) { dto.TemplateName = gofakeit.ProductName() })

	// Act
	state.AutoRegenerate(now)

	// Assert
	handlerMock.AssertNumberOfCalls(t, "GenerateTemplate", 1)
}

func TestWhenDebounceHasNotElapsed_KeepsWaiting(t *testing.T) {
	t.Parallel()
	// Arrange
	state, _ := newAutoRegenerateState()
	now := time.Now()
	state.AutoRegenerate(now)
	state.UpdateState(func(dto *editor_state_model.EditorState) { dto.TemplateName = gofakeit.ProductName() })
	state.AutoRegenerate(now)

	// Act
	_, scheduleRedraw := state.AutoRegenerate(now.Add(100 * time.Millisecond))

	// Assert
	assert.True(t, scheduleRedraw)
}

func TestWhenDebounceElapsesWithoutFurtherEdits_Regenerates(t *testing.T) {
	t.Parallel()
	// Arrange
	state, handlerMock := newAutoRegenerateState()
	now := time.Now()
	state.AutoRegenerate(now)
	state.UpdateState(func(dto *editor_state_model.EditorState) { dto.TemplateName = gofakeit.ProductName() })
	state.AutoRegenerate(now)

	// Act
	state.AutoRegenerate(now.Add(time.Second))

	// Assert
	handlerMock.AssertNumberOfCalls(t, "GenerateTemplate", 2)
}

func TestWhenManualEditsAreStored_RegenerationDecisionReceivesNoManualEdits(t *testing.T) {
	t.Parallel()
	// Arrange
	recorder := &decisionRecorder{IRegenerationHandler: test_helpers.NewRegenerationHandler()}
	state, _ := newAutoRegenerateStateWith(recorder)
	now := time.Now()
	storeManualZone(state)

	// Act
	driveIdleAndDebouncedFrames(state, now)

	// Assert
	assert.Equal(t,
		decisionObservation{SawPending: true, SawManualEdits: false},
		recorder.observe())
}

func TestWhenManualEditsAreStored_TheySurviveIdleAndDebouncedFrames(t *testing.T) {
	t.Parallel()
	// Arrange
	state, _ := newAutoRegenerateState()
	now := time.Now()
	storeManualZone(state)

	// Act
	driveIdleAndDebouncedFrames(state, now)

	// Assert
	currentState := state.GetStateData()
	assert.True(t, currentState.HasManualEdits())
}

// newAutoRegenerateState returns a State wired to a mock whose GenerateTemplate
// always succeeds, plus the mock for call-count assertions.
func newAutoRegenerateState() (*drivers.State, *test_helpers.TemplateHandlerMock) {
	return newAutoRegenerateStateWith(test_helpers.NewRegenerationHandler())
}

func newAutoRegenerateStateWith(
	regeneration handler_interfaces.IRegenerationHandler) (*drivers.State, *test_helpers.TemplateHandlerMock) {
	handlerMock := &test_helpers.TemplateHandlerMock{}
	template := test_helpers.GetDefaultTemplateModel()
	handlerMock.On("GenerateTemplate", mock.Anything).Return(dtos.TemplateLoadDto{Template: &template}, nil)
	handlerMock.On("UpdateTemplate", mock.Anything).Return(dtos.TemplateLoadDto{Template: &template}, nil)
	return drivers.NewUIState(
		handlerMock,
		test_helpers.NewFileSystemHandler(),
		regeneration,
		false), handlerMock
}

func storeManualZone(state *drivers.State) {
	state.UpdateState(func(current *editor_state_model.EditorState) {
		current.ManualZones = []template_model.Zone{{Name: "Neutral-A"}}
	})
}

// driveIdleAndDebouncedFrames generates once, runs an idle frame, then edits a
// debounced option and runs two more frames so a pending next state exists.
func driveIdleAndDebouncedFrames(state *drivers.State, now time.Time) {
	state.AutoRegenerate(now)
	state.AutoRegenerate(now)
	state.UpdateState(func(current *editor_state_model.EditorState) { current.TemplateName = gofakeit.ProductName() })
	state.AutoRegenerate(now)
	state.AutoRegenerate(now)
}

type decisionObservation struct {
	SawPending     bool
	SawManualEdits bool
}

// decisionRecorder forwards to the real regeneration handler and keeps every
// per-frame decision request it was given.
type decisionRecorder struct {
	handler_interfaces.IRegenerationHandler

	requests []dtos.RegenerationDecisionRequestDto
}

func (this *decisionRecorder) DecideRegeneration(
	request dtos.RegenerationDecisionRequestDto) dtos.RegenerationDecisionDto {
	this.requests = append(this.requests, request)
	return this.IRegenerationHandler.DecideRegeneration(request)
}

func (this *decisionRecorder) observe() decisionObservation {
	var observation decisionObservation
	for _, request := range this.requests {
		observation.SawPending = observation.SawPending || request.Next != nil
		for _, snapshot := range []*editor_state_dto.EditorStateDto{request.Previous, request.Current, request.Next} {
			if snapshot != nil && snapshot.HasManualEdits() {
				observation.SawManualEdits = true
			}
		}
	}
	return observation
}
