package stateHandler_test

import (
	"errors"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_errors"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestWhenStatePathIsEmpty_ReturnsNoOutputPathError(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := handlers.NewStateHandler(&test_helpers.FileServiceMock{}, newPassingValidator())

	// Act
	_, _, err := handler.LoadState("", true)

	// Assert
	assert.ErrorIs(t, err, common_errors.ErrNoOutputPath)
}

func TestWhenStatePathIsWhitespaceOnly_ReturnsNoOutputPathError(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := handlers.NewStateHandler(&test_helpers.FileServiceMock{}, newPassingValidator())

	// Act
	_, _, err := handler.LoadState("  \t ", true)

	// Assert
	assert.ErrorIs(t, err, common_errors.ErrNoOutputPath)
}

func TestWhenStatePathIsPadded_LoadsTheTrimmedPath(t *testing.T) {
	t.Parallel()
	// Arrange
	path := gofakeit.Word() + ".gen.json"
	state := dtos.NewDefaultEditorStateDto()
	fileService := &test_helpers.FileServiceMock{}
	fileService.On("LoadSettingsFile", path).Return(&state, nil)
	handler := handlers.NewStateHandler(fileService, newPassingValidator())

	// Act
	_, _, _ = handler.LoadState("  "+path+"  ", true)

	// Assert
	fileService.AssertCalled(t, "LoadSettingsFile", path)
}

func TestWhenSettingsFileCannotBeLoaded_PropagatesTheError(t *testing.T) {
	t.Parallel()
	// Arrange
	expectedError := errors.New(gofakeit.Sentence(3))
	fileService := &test_helpers.FileServiceMock{}
	fileService.On("LoadSettingsFile", mock.Anything).Return(nil, expectedError)
	handler := handlers.NewStateHandler(fileService, newPassingValidator())

	// Act
	_, _, err := handler.LoadState(gofakeit.Word(), true)

	// Assert
	assert.ErrorIs(t, err, expectedError)
}

func TestWhenSettingsFileCannotBeLoaded_ReturnsNoState(t *testing.T) {
	t.Parallel()
	// Arrange
	fileService := &test_helpers.FileServiceMock{}
	fileService.On("LoadSettingsFile", mock.Anything).Return(nil, errors.New(gofakeit.Sentence(3)))
	handler := handlers.NewStateHandler(fileService, newPassingValidator())

	// Act
	state, _, _ := handler.LoadState(gofakeit.Word(), true)

	// Assert
	assert.Nil(t, state)
}

func TestWhenSettingsFileIsLoaded_ReturnsTheValidatedState(t *testing.T) {
	t.Parallel()
	// Arrange
	loaded := dtos.NewDefaultEditorStateDto()
	loaded.TemplateName = gofakeit.Word()
	fileService := &test_helpers.FileServiceMock{}
	fileService.On("LoadSettingsFile", mock.Anything).Return(&loaded, nil)
	handler := handlers.NewStateHandler(fileService, newPassingValidator())

	// Act
	state, _, err := handler.LoadState(gofakeit.Word(), true)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, loaded, *state)
}

func TestWhenValidationReportsIssues_ReturnsThemAsWarnings(t *testing.T) {
	t.Parallel()
	// Arrange
	firstMessage := gofakeit.Sentence(3)
	secondMessage := gofakeit.Sentence(3)
	loaded := dtos.NewDefaultEditorStateDto()
	fileService := &test_helpers.FileServiceMock{}
	fileService.On("LoadSettingsFile", mock.Anything).Return(&loaded, nil)
	handler := handlers.NewStateHandler(fileService, newValidatorReporting(firstMessage, secondMessage))

	// Act
	_, warnings, _ := handler.LoadState(gofakeit.Word(), false)

	// Assert
	assert.Equal(t, []string{firstMessage, secondMessage}, warnings)
}
