package state_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/drivers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
)

func TestWhenStateIsCreated_StateDataIsDefault(t *testing.T) {
	t.Parallel()
	// Arrange
	handlerMock := &test_helpers.TemplateHandlerMock{}

	// Act
	state := drivers.NewUIState(
		handlerMock, test_helpers.NewFileSystemHandler(), test_helpers.NewRegenerationHandler(), false)

	// Assert
	assert.Equal(t, editor_state_model.NewDefaultEditorStateModel(), state.GetStateData())
}

func TestWhenStateIsCreated_NoDialogIsOpen(t *testing.T) {
	t.Parallel()
	// Arrange
	handlerMock := &test_helpers.TemplateHandlerMock{}

	// Act
	state := drivers.NewUIState(
		handlerMock, test_helpers.NewFileSystemHandler(), test_helpers.NewRegenerationHandler(), false)

	// Assert
	assert.False(t, state.GetDialogHost().IsOpen())
}

func TestWhenTemplateDirLookupIsSkipped_OutputPathIsEmpty(t *testing.T) {
	t.Parallel()
	// Arrange
	handlerMock := &test_helpers.TemplateHandlerMock{}

	// Act
	state := drivers.NewUIState(
		handlerMock, test_helpers.NewFileSystemHandler(), test_helpers.NewRegenerationHandler(), false)

	// Assert
	assert.Empty(t, state.GetOutputPath())
}

// The lookup falls back to the working directory, so the path is set even on a
// machine without the game installed.
func TestWhenTemplateDirLookupIsRequested_OutputPathIsSet(t *testing.T) {
	t.Parallel()
	// Arrange
	handlerMock := &test_helpers.TemplateHandlerMock{}

	// Act
	state := drivers.NewUIState(
		handlerMock, test_helpers.NewFileSystemHandler(), test_helpers.NewRegenerationHandler(), true)

	// Assert
	assert.NotEmpty(t, state.GetOutputPath())
}
