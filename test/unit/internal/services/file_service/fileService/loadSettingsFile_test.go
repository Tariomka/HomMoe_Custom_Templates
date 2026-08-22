package fileService_test

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenSettingsFileIsRequested_LoadsItFromTheGivenPath(t *testing.T) {
	t.Parallel()
	// Arrange
	service, mocks := newServiceWithMocks()
	settingsPath := filepath.Join("any", "where", "state.gen.json")
	mocks.editorState.On("Load", settingsPath).Return(editor_state_dto.NewDefaultEditorStateDto(), nil)

	// Act
	_, err := service.LoadSettingsFile(settingsPath)

	// Assert
	require.NoError(t, err)
	mocks.editorState.AssertCalled(t, "Load", settingsPath)
}

func TestWhenSettingsFileIsLoaded_ReturnsTheLoadedState(t *testing.T) {
	t.Parallel()
	// Arrange
	service, mocks := newServiceWithMocks()
	expected := editor_state_model.NewDefaultEditorStateModel()
	expected.TemplateName = "Loaded"
	mocks.editorState.On("Load", "state.gen.json").
		Return(editor_state_dto.NewEditorStateDto(expected), nil)

	// Act
	actual, err := service.LoadSettingsFile("state.gen.json")

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expected, *actual)
}

func TestWhenSettingsFileCannotBeLoaded_ReturnsError(t *testing.T) {
	t.Parallel()
	// Arrange
	service, mocks := newServiceWithMocks()
	expectedError := errors.New("unreadable")
	mocks.editorState.On("Load", "state.gen.json").Return(editor_state_dto.EditorStateDto{}, expectedError)

	// Act
	_, err := service.LoadSettingsFile("state.gen.json")

	// Assert
	assert.ErrorIs(t, err, expectedError)
}
