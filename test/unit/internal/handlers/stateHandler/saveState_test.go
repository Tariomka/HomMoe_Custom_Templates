package stateHandler_test

import (
	"errors"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_errors"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestWhenThereIsNoStateToSave_ReturnsNothingToSaveError(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := handlers.NewStateHandler(&test_helpers.FileServiceMock{}, newPassingValidator())

	// Act
	_, err := handler.SaveState(editor_state_dto.EditorStateSaveDto{OutputPath: gofakeit.Word()})

	// Assert
	assert.ErrorIs(t, err, common_errors.ErrNothingToSave)
}

func TestWhenStateOutputPathIsEmpty_ReturnsNoOutputPathError(t *testing.T) {
	t.Parallel()
	// Arrange
	state := editor_state_model.NewDefaultEditorStateModel()
	handler := handlers.NewStateHandler(&test_helpers.FileServiceMock{}, newPassingValidator())

	// Act
	_, err := handler.SaveState(editor_state_dto.EditorStateSaveDto{State: &state, OutputPath: ""})

	// Assert
	assert.ErrorIs(t, err, common_errors.ErrNoOutputPath)
}

func TestWhenStateOutputPathIsWhitespaceOnly_ReturnsNoOutputPathError(t *testing.T) {
	t.Parallel()
	// Arrange
	state := editor_state_model.NewDefaultEditorStateModel()
	handler := handlers.NewStateHandler(&test_helpers.FileServiceMock{}, newPassingValidator())

	// Act
	_, err := handler.SaveState(editor_state_dto.EditorStateSaveDto{State: &state, OutputPath: " \t "})

	// Assert
	assert.ErrorIs(t, err, common_errors.ErrNoOutputPath)
}

func TestWhenStateOutputPathIsPadded_SavesToTheTrimmedPath(t *testing.T) {
	t.Parallel()
	// Arrange
	outputPath := gofakeit.Word()
	state := editor_state_model.NewDefaultEditorStateModel()
	fileService := &test_helpers.FileServiceMock{}
	fileService.On("SaveSettings", outputPath, &state).Return(gofakeit.Word(), nil)
	handler := handlers.NewStateHandler(fileService, newPassingValidator())

	// Act
	_, _ = handler.SaveState(editor_state_dto.EditorStateSaveDto{State: &state, OutputPath: " " + outputPath + " "})

	// Assert
	fileService.AssertCalled(t, "SaveSettings", outputPath, &state)
}

func TestWhenStateIsSaved_ReturnsTheWrittenPath(t *testing.T) {
	t.Parallel()
	// Arrange
	expectedPath := gofakeit.Word() + ".gen.json"
	state := editor_state_model.NewDefaultEditorStateModel()
	fileService := &test_helpers.FileServiceMock{}
	fileService.On("SaveSettings", mock.Anything, mock.Anything).Return(expectedPath, nil)
	handler := handlers.NewStateHandler(fileService, newPassingValidator())

	// Act
	writtenPath, err := handler.SaveState(
		editor_state_dto.EditorStateSaveDto{State: &state, OutputPath: gofakeit.Word()},
	)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedPath, writtenPath)
}

func TestWhenStateCannotBeSaved_PropagatesTheError(t *testing.T) {
	t.Parallel()
	// Arrange
	expectedError := errors.New(gofakeit.Sentence(3))
	state := editor_state_model.NewDefaultEditorStateModel()
	fileService := &test_helpers.FileServiceMock{}
	fileService.On("SaveSettings", mock.Anything, mock.Anything).Return("", expectedError)
	handler := handlers.NewStateHandler(fileService, newPassingValidator())

	// Act
	_, err := handler.SaveState(editor_state_dto.EditorStateSaveDto{State: &state, OutputPath: gofakeit.Word()})

	// Assert
	assert.ErrorIs(t, err, expectedError)
}
