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
	mocks.migrator.On("Load", settingsPath).Return(editor_state.EditorState{}, nil)

	// Act
	_, err := service.LoadSettingsFile(settingsPath)

	// Assert
	require.NoError(t, err)
	mocks.migrator.AssertCalled(t, "Load", settingsPath)
}

// The migrator decodes into whatever it is handed, so seeding the defaults here
// is what lets a key the file omits keep its default rather than collapse to a
// zero value.
func TestWhenSettingsFileIsRequested_TheDecodeIsSeededWithTheDefaultEntity(t *testing.T) {
	t.Parallel()
	// Arrange
	service, mocks := newServiceWithMocks()
	mocks.migrator.On("Load", "state.gen.json").Return(editor_state.EditorState{}, nil)

	// Act
	_, err := service.LoadSettingsFile("state.gen.json")

	// Assert
	require.NoError(t, err)
	assert.Equal(t, mocks.mapper.NewDefaultEntity(), mocks.migrator.seed)
}

func TestWhenSettingsFileIsLoaded_ReturnsTheLoadedStateAsAModel(t *testing.T) {
	t.Parallel()
	// Arrange
	service, mocks := newServiceWithMocks()
	loaded := mocks.mapper.NewDefaultEntity()
	loaded.TemplateName = "Loaded"
	mocks.migrator.On("Load", "state.gen.json").Return(loaded, nil)

	// Act
	actual, err := service.LoadSettingsFile("state.gen.json")

	// Assert
	require.NoError(t, err)
	assert.Equal(t, mocks.mapper.ToModel(loaded), *actual)
}

// The migrator is the only thing that can fail a load - a missing file, an
// unreadable one and a file too new for this build all arrive as its error.
func TestWhenSettingsFileCannotBeLoaded_ReturnsError(t *testing.T) {
	t.Parallel()
	// Arrange
	service, mocks := newServiceWithMocks()
	expectedError := errors.New("unreadable")
	mocks.migrator.On("Load", "state.gen.json").Return(editor_state.EditorState{}, expectedError)

	// Act
	_, err := service.LoadSettingsFile("state.gen.json")

	// Assert
	assert.ErrorIs(t, err, expectedError)
}
