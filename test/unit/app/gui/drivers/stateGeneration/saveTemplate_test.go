package stateGeneration_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/drivers"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWhenTemplateSaveSucceeds_StatusReportsSavedPath(t *testing.T) {
	t.Parallel()
	// Arrange
	handlerMock := &test_helpers.TemplateHandlerMock{}
	savedPath := gofakeit.Word() + ".rmg.json"
	handlerMock.On("SaveTemplate", mock.Anything).Return(savedPath, nil)
	state := drivers.NewUIState(handlerMock, test_helpers.NewFileSystemHandler(), false)

	// Act
	state.SaveTemplate()

	// Assert
	message, isError := state.GetStatus()
	assert.Equal(t, []any{"Saved template to " + savedPath, false}, []any{message, isError})
}

func TestWhenTemplateSaveFailsBeforeWriting_ErrorStatusIsSet(t *testing.T) {
	t.Parallel()
	// Arrange
	handlerMock := &test_helpers.TemplateHandlerMock{}
	handlerMock.On("SaveTemplate", mock.Anything).Return("", gofakeit.ErrorValidation())
	state := drivers.NewUIState(handlerMock, test_helpers.NewFileSystemHandler(), false)

	// Act
	state.SaveTemplate()

	// Assert
	message, isError := state.GetStatus()
	assert.True(t, isError, "status: %s", message)
}

func TestWhenTemplateSavedButPreviewFails_StatusWarnsAboutPreview(t *testing.T) {
	t.Parallel()
	// Arrange
	handlerMock := &test_helpers.TemplateHandlerMock{}
	savedPath := gofakeit.Word() + ".rmg.json"
	handlerMock.On("SaveTemplate", mock.Anything).Return(savedPath, gofakeit.ErrorValidation())
	state := drivers.NewUIState(handlerMock, test_helpers.NewFileSystemHandler(), false)

	// Act
	state.SaveTemplate()

	// Assert
	message, _ := state.GetStatus()
	assert.Contains(t, message, "failed to write preview PNG")
}
