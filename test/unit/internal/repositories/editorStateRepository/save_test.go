package editorStateRepository_test

import (
	"math"
	"os"
	"path/filepath"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/repositories"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newStateWithNaN() editor_state_dto.EditorStateDto {
	state := editor_state_dto.NewDefaultEditorStateDto()
	state.PlayerZoneSize = math.NaN()

	return state
}

func TestWhenStateIsSaved_ReturnsPathWithGenJsonExtension(t *testing.T) {
	t.Parallel()
	// Arrange
	outputDir := t.TempDir()

	// Act
	writtenPath, err := repositories.NewEditorStateRepository().Save(
		outputDir, "My_State", editor_state_dto.NewDefaultEditorStateDto())

	// Assert
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(outputDir, "My_State.gen.json"), writtenPath)
}

func TestWhenStateIsSaved_WritesIndentedJson(t *testing.T) {
	t.Parallel()
	// Arrange
	outputDir := t.TempDir()
	writtenPath, err := repositories.NewEditorStateRepository().Save(
		outputDir, "State", editor_state_dto.NewDefaultEditorStateDto())
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
	writtenPath, err := repositories.NewEditorStateRepository().Save(
		outputDir, "a/b:c", editor_state_dto.NewDefaultEditorStateDto())

	// Assert
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(outputDir, "a_b_c.gen.json"), writtenPath)
}

func TestWhenStateNameIsOnlyWhitespace_FallsBackToGeneratedTemplateFileName(t *testing.T) {
	t.Parallel()
	// Arrange
	outputDir := t.TempDir()

	// Act
	writtenPath, err := repositories.NewEditorStateRepository().Save(
		outputDir, "   ", editor_state_dto.NewDefaultEditorStateDto())

	// Assert
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(outputDir, "Generated_Template.gen.json"), writtenPath)
}

func TestWhenStateDirectoryIsMissing_CreatesIt(t *testing.T) {
	t.Parallel()
	// Arrange
	outputDir := filepath.Join(t.TempDir(), "nested", "state")

	// Act
	_, err := repositories.NewEditorStateRepository().
		Save(outputDir, "State", editor_state_dto.NewDefaultEditorStateDto())

	// Assert
	require.NoError(t, err)
	assert.DirExists(t, outputDir)
}

func TestWhenSavedStateIsLoaded_RoundTripsState(t *testing.T) {
	t.Parallel()
	// Arrange
	repository := repositories.NewEditorStateRepository()
	state := editor_state_dto.NewDefaultEditorStateDto()
	state.TemplateName = gofakeit.ProductName()
	state.PlayerCount = gofakeit.Number(2, 8)
	writtenPath, err := repository.Save(t.TempDir(), "State", state)
	require.NoError(t, err)

	// Act
	loaded, loadErr := repository.Load(writtenPath)

	// Assert
	require.NoError(t, loadErr)
	assert.Equal(t, state, loaded)
}

func TestWhenStateContainsNaNValue_ReturnsError(t *testing.T) {
	t.Parallel()
	// Arrange
	outputDir := t.TempDir()

	// Act
	_, err := repositories.NewEditorStateRepository().Save(outputDir, "State", newStateWithNaN())

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
	_, err := repositories.NewEditorStateRepository().Save(outputDir, "State", newStateWithNaN())

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
	_, err := repositories.NewEditorStateRepository().Save(outputDir, "State", newStateWithNaN())

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
	_, err := repositories.NewEditorStateRepository().
		Save(outputDir, "State", editor_state_dto.NewDefaultEditorStateDto())

	// Assert
	assert.Error(t, err)
}
