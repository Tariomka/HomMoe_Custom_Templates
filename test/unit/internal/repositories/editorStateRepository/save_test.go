package editorStateRepository_test

import (
	"math"
	"os"
	"path/filepath"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenStateIsSaved_ReturnsPathWithGenJsonExtension(t *testing.T) {
	t.Parallel()
	// Arrange
	outputDir := t.TempDir()

	// Act
	writtenPath, err := newRepository().Save(outputDir, "My_State", editor_state.EditorState{})

	// Assert
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(outputDir, "My_State.gen.json"), writtenPath)
}

func TestWhenStateIsSaved_WritesIndentedJson(t *testing.T) {
	t.Parallel()
	// Arrange
	outputDir := t.TempDir()
	writtenPath, err := newRepository().Save(outputDir, "State", editor_state.EditorState{})
	require.NoError(t, err)

	// Act
	data, readErr := os.ReadFile(writtenPath)

	// Assert
	require.NoError(t, readErr)
	assert.Contains(t, string(data), "\n  ")
}

func TestWhenStateNameContainsInvalidCharacters_WritesUnderASanitizedName(t *testing.T) {
	t.Parallel()
	// Arrange
	outputDir := t.TempDir()

	// Act
	writtenPath, err := newRepository().Save(outputDir, "a/b:c", editor_state.EditorState{})

	// Assert
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(outputDir, "a_b_c.gen.json"), writtenPath)
}

func TestWhenStateNameIsOnlyWhitespace_FallsBackToGeneratedTemplateFileName(t *testing.T) {
	t.Parallel()
	// Arrange
	outputDir := t.TempDir()

	// Act
	writtenPath, err := newRepository().Save(outputDir, "   ", editor_state.EditorState{})

	// Assert
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(outputDir, "Generated_Template.gen.json"), writtenPath)
}

func TestWhenStateDirectoryIsMissing_CreatesIt(t *testing.T) {
	t.Parallel()
	// Arrange
	outputDir := filepath.Join(t.TempDir(), "nested", "state")

	// Act
	_, err := newRepository().Save(outputDir, "State", editor_state.EditorState{})

	// Assert
	require.NoError(t, err)
	assert.DirExists(t, outputDir)
}

func TestWhenSavedStateIsLoaded_RoundTripsState(t *testing.T) {
	t.Parallel()
	// Arrange
	repository := newRepository()
	state := editor_state.EditorState{
		SchemaVersion: editor_state.CurrentEditorStateSchemaVersion,
		TemplateName:  gofakeit.ProductName(),
		PlayerCount:   gofakeit.Number(2, 8),
	}
	writtenPath, err := repository.Save(t.TempDir(), "State", state)
	require.NoError(t, err)
	loaded := editor_state.EditorState{}

	// Act
	loadErr := repository.Load(writtenPath, &loaded)

	// Assert
	require.NoError(t, loadErr)
	assert.Equal(t, state, loaded)
}

func TestWhenStateContainsNaNValue_ReturnsError(t *testing.T) {
	t.Parallel()
	// Arrange
	outputDir := t.TempDir()

	// Act
	_, err := newRepository().Save(outputDir, "State", newStateWithNaN())

	// Assert
	assert.Error(t, err)
}

func TestWhenEncodingFailsOverAnExistingState_LeavesTheDestinationUntouched(t *testing.T) {
	t.Parallel()
	// Arrange
	outputDir := t.TempDir()
	destinationPath := filepath.Join(outputDir, "State.gen.json")
	require.NoError(t, os.WriteFile(destinationPath, []byte("previous"), 0o644))

	// Act
	_, err := newRepository().Save(outputDir, "State", newStateWithNaN())

	// Assert
	require.Error(t, err)
	survivingContent, readErr := os.ReadFile(destinationPath)
	require.NoError(t, readErr)
	assert.Equal(t, "previous", string(survivingContent))
}

func TestWhenStateEncodingFails_RemovesTheTemporaryFile(t *testing.T) {
	t.Parallel()
	// Arrange
	outputDir := t.TempDir()

	// Act
	_, err := newRepository().Save(outputDir, "State", newStateWithNaN())

	// Assert
	require.Error(t, err)
	assert.NoFileExists(t, filepath.Join(outputDir, "TEMP-State.gen.json"))
}

func TestWhenStateParentPathIsAFile_ReturnsError(t *testing.T) {
	t.Parallel()
	// Arrange
	blockerPath := filepath.Join(t.TempDir(), "blocker")
	require.NoError(t, os.WriteFile(blockerPath, []byte("x"), 0o644))
	outputDir := filepath.Join(blockerPath, "child")

	// Act
	_, err := newRepository().Save(outputDir, "State", editor_state.EditorState{})

	// Assert
	assert.Error(t, err)
}

func newStateWithNaN() editor_state.EditorState {
	return editor_state.EditorState{PlayerZoneSize: math.NaN()}
}
