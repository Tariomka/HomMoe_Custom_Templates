package stateGeneration_test

import (
	"testing"
	"time"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/drivers"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
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
	state.UpdateState(func(dto *editor_state_dto.EditorStateDto) { dto.PlayerCount++ })

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
	state.UpdateState(func(dto *editor_state_dto.EditorStateDto) { dto.TemplateName = gofakeit.ProductName() })

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
	state.UpdateState(func(dto *editor_state_dto.EditorStateDto) { dto.TemplateName = gofakeit.ProductName() })

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
	state.UpdateState(func(dto *editor_state_dto.EditorStateDto) { dto.TemplateName = gofakeit.ProductName() })
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
	state.UpdateState(func(dto *editor_state_dto.EditorStateDto) { dto.TemplateName = gofakeit.ProductName() })
	state.AutoRegenerate(now)

	// Act
	state.AutoRegenerate(now.Add(time.Second))

	// Assert
	handlerMock.AssertNumberOfCalls(t, "GenerateTemplate", 2)
}

// newAutoRegenerateState returns a State wired to a mock whose GenerateTemplate
// always succeeds, plus the mock for call-count assertions.
func newAutoRegenerateState() (*drivers.State, *test_helpers.TemplateHandlerMock) {
	handlerMock := &test_helpers.TemplateHandlerMock{}
	template := test_helpers.GetDefaultTemplate()
	handlerMock.On("GenerateTemplate", mock.Anything).Return(dtos.TemplateLoadDto{Template: &template}, nil)
	return drivers.NewUIState(
		handlerMock, test_helpers.NewFileSystemHandler(), test_helpers.NewRegenerationHandler(), false), handlerMock
}
