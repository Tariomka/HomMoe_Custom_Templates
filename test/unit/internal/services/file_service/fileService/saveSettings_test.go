package fileService_test

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenSettingsAreSaved_UsesTheDirectoryOfTheGivenPath(t *testing.T) {
	t.Parallel()
	// Arrange
	service, mocks := newServiceWithMocks()
	outputDirectory := filepath.Join("out", "states")
	state := dtos.NewDefaultEditorStateDto()
	state.TemplateName = "My Template"
	mocks.editorState.On("Save", outputDirectory, "My Template", state).Return("written", nil)

	// Act
	_, err := service.SaveSettings(filepath.Join(outputDirectory, "ignored.gen.json"), &state)

	// Assert
	require.NoError(t, err)
	mocks.editorState.AssertCalled(t, "Save", outputDirectory, "My Template", state)
}

func TestWhenStateTemplateNameNeedsSanitizing_ForwardsItUnchangedToTheRepository(t *testing.T) {
	t.Parallel()
	// Arrange
	service, mocks := newServiceWithMocks()
	state := dtos.NewDefaultEditorStateDto()
	state.TemplateName = "a/b:c"
	mocks.editorState.On("Save", ".", "a/b:c", state).Return("written", nil)

	// Act
	_, err := service.SaveSettings("ignored.gen.json", &state)

	// Assert
	require.NoError(t, err)
	mocks.editorState.AssertCalled(t, "Save", ".", "a/b:c", state)
}

func TestWhenSettingsAreSaved_ReturnsThePathTheRepositoryWrote(t *testing.T) {
	t.Parallel()
	// Arrange
	service, mocks := newServiceWithMocks()
	state := dtos.NewDefaultEditorStateDto()
	state.TemplateName = "Name"
	expectedPath := filepath.Join("out", "Name.gen.json")
	mocks.editorState.On("Save", "out", "Name", state).Return(expectedPath, nil)

	// Act
	actualPath, err := service.SaveSettings(filepath.Join("out", "whatever.gen.json"), &state)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedPath, actualPath)
}

func TestWhenSettingsCannotBeSaved_ReturnsError(t *testing.T) {
	t.Parallel()
	// Arrange
	service, mocks := newServiceWithMocks()
	state := dtos.NewDefaultEditorStateDto()
	state.TemplateName = "Name"
	expectedError := errors.New("disk full")
	mocks.editorState.On("Save", "out", "Name", state).Return("", expectedError)

	// Act
	_, err := service.SaveSettings(filepath.Join("out", "whatever.gen.json"), &state)

	// Assert
	assert.ErrorIs(t, err, expectedError)
}
