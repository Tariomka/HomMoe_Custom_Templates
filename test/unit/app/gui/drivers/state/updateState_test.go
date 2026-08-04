package state_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/drivers"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// newGeneratedState returns a State that has generated once, so a previous
// state snapshot exists and change detection is active.
func newGeneratedState() *drivers.State {
	handlerMock := &test_helpers.TemplateHandlerMock{}
	template := test_helpers.GetDefaultTemplate()
	handlerMock.On("GenerateTemplate", mock.Anything).Return(dtos.TemplateLoadDto{Template: &template}, nil)
	state := drivers.NewUIState(handlerMock, false)
	state.Generate()
	return state
}

func TestWhenUpdateChangesState_ChangeIsApplied(t *testing.T) {
	t.Parallel()
	// Arrange
	state := drivers.NewUIState(&test_helpers.TemplateHandlerMock{}, false)
	playerCount := gofakeit.Number(3, 8)

	// Act
	state.UpdateState(func(dto *dtos.EditorStateDto) { dto.PlayerCount = playerCount })

	// Assert
	assert.Equal(t, playerCount, state.GetStateData().PlayerCount)
}

func TestWhenUpdateChangesStateAfterGeneration_StateBecomesUnsaved(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newGeneratedState()

	// Act
	state.UpdateState(func(dto *dtos.EditorStateDto) { dto.TemplateName = gofakeit.ProductName() })

	// Assert
	assert.True(t, state.IsUnsaved())
}

func TestWhenUpdateDoesNotChangeStateAfterGeneration_StateStaysSaved(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newGeneratedState()

	// Act
	state.UpdateState(func(_ *dtos.EditorStateDto) {})

	// Assert
	assert.False(t, state.IsUnsaved())
}
