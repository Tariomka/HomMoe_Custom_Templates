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

// newDirtyGeneratedState returns a State with a generated template and an
// unsaved change, so Reset has something to clear in every dimension.
func newDirtyGeneratedState() *drivers.State {
	handlerMock := &test_helpers.TemplateHandlerMock{}
	template := test_helpers.GetDefaultTemplate()
	handlerMock.On("GenerateTemplate", mock.Anything).Return(dtos.TemplateLoadDto{Template: &template}, nil)
	state := drivers.NewUIStateWithHandler(handlerMock)
	state.Generate()
	state.UpdateState(func(dto *dtos.EditorStateDto) { dto.TemplateName = gofakeit.ProductName() })
	return state
}

func TestWhenStateIsReset_UnsavedFlagIsCleared(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newDirtyGeneratedState()

	// Act
	state.Reset()

	// Assert
	assert.False(t, state.IsUnsaved())
}

func TestWhenStateIsReset_LastTemplateIsCleared(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newDirtyGeneratedState()

	// Act
	state.Reset()

	// Assert
	assert.Nil(t, state.GetLastTemplate())
}

func TestWhenStateIsReset_StateDataReturnsToDefault(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newDirtyGeneratedState()

	// Act
	state.Reset()

	// Assert
	assert.Equal(t, dtos.NewDefaultEditorStateDto(), state.GetStateData())
}

func TestWhenStateIsReset_StatusReportsNewFile(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newDirtyGeneratedState()

	// Act
	state.Reset()

	// Assert
	message, isError := state.GetStatus()
	assert.Equal(t, []any{"New settings file.", false}, []any{message, isError})
}
