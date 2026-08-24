package fileService_test

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenSettingsFileIsRequested_LoadsItFromTheGivenPath(t *testing.T) {
	t.Parallel()
	// Arrange
	service, mocks := newServiceWithMocks()
	settingsPath := filepath.Join("any", "where", "state.gen.json")
	mocks.editorState.On("Load", settingsPath).Return(editor_state.EditorState{}, nil)

	// Act
	_, err := service.LoadSettingsFile(settingsPath)

	// Assert
	require.NoError(t, err)
	mocks.editorState.AssertCalled(t, "Load", settingsPath)
}

// The repository decodes into whatever it is handed, so seeding the defaults
// here is what lets a key the file omits keep its default rather than collapse
// to a zero value.
func TestWhenSettingsFileIsRequested_TheDecodeIsSeededWithTheDefaultEntity(t *testing.T) {
	t.Parallel()
	// Arrange
	service, mocks := newServiceWithMocks()
	mocks.editorState.On("Load", "state.gen.json").Return(editor_state.EditorState{}, nil)

	// Act
	_, err := service.LoadSettingsFile("state.gen.json")

	// Assert
	require.NoError(t, err)
	assert.Equal(t, mocks.mapper.NewDefaultEntity(), mocks.editorState.seed)
}

func TestWhenSettingsFileIsLoaded_ReturnsTheLoadedStateAsAModel(t *testing.T) {
	t.Parallel()
	// Arrange
	service, mocks := newServiceWithMocks()
	loaded := mocks.mapper.NewDefaultEntity()
	loaded.TemplateName = "Loaded"
	mocks.editorState.On("Load", "state.gen.json").Return(loaded, nil)

	// Act
	actual, err := service.LoadSettingsFile("state.gen.json")

	// Assert
	require.NoError(t, err)
	assert.Equal(t, mocks.mapper.ToModel(loaded), *actual)
}

func TestWhenSettingsFileCannotBeLoaded_ReturnsError(t *testing.T) {
	t.Parallel()
	// Arrange
	service, mocks := newServiceWithMocks()
	expectedError := errors.New("unreadable")
	mocks.editorState.On("Load", "state.gen.json").Return(editor_state.EditorState{}, expectedError)

	// Act
	_, err := service.LoadSettingsFile("state.gen.json")

	// Assert
	assert.ErrorIs(t, err, expectedError)
}
