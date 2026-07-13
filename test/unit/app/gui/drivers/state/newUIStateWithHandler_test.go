package state_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/drivers"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
)

func TestWhenStateIsCreated_StateDataIsDefault(t *testing.T) {
	t.Parallel()
	// Arrange
	handlerMock := &test_helpers.TemplateHandlerMock{}

	// Act
	state := drivers.NewUIStateWithHandler(handlerMock)

	// Assert
	assert.Equal(t, dtos.NewDefaultEditorStateDto(), state.GetStateData())
}

func TestWhenStateIsCreated_NoDialogIsOpen(t *testing.T) {
	t.Parallel()
	// Arrange
	handlerMock := &test_helpers.TemplateHandlerMock{}

	// Act
	state := drivers.NewUIStateWithHandler(handlerMock)

	// Assert
	assert.False(t, state.GetDialogHost().IsOpen())
}

func TestWhenStateIsCreated_OutputPathIsEmpty(t *testing.T) {
	t.Parallel()
	// Arrange
	handlerMock := &test_helpers.TemplateHandlerMock{}

	// Act
	state := drivers.NewUIStateWithHandler(handlerMock)

	// Assert
	assert.Empty(t, state.GetOutputPath())
}
