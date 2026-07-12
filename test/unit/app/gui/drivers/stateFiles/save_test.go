package stateFiles_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/drivers"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
)

func TestWhenNoCurrentPathExists_SaveOpensSaveAsDialog(t *testing.T) {
	t.Parallel()
	// Arrange
	handlerMock := &test_helpers.TemplateHandlerMock{}
	state := drivers.NewUIStateWithHandler(handlerMock)

	// Act
	state.Save()

	// Assert
	assert.True(t, state.GetDialogHost().IsOpen())
}

func TestWhenNoCurrentPathExists_SaveDoesNotCallHandler(t *testing.T) {
	t.Parallel()
	// Arrange
	handlerMock := &test_helpers.TemplateHandlerMock{}
	state := drivers.NewUIStateWithHandler(handlerMock)

	// Act
	state.Save()

	// Assert
	handlerMock.AssertNotCalled(t, "SaveState")
}
