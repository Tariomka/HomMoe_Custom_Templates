package guiHandler_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/common"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenStateFilePathIsEmpty_ReturnsNoOutputPathError(t *testing.T) {
	// Arrange
	handler := handlers.NewGuiHandler()

	// Act
	_, err := handler.LoadState("")

	// Assert
	assert.ErrorIs(t, err, common.ErrNoOutputPath)
}

func TestWhenStateFilePathIsWhitespaceOnly_ReturnsNoOutputPathError(t *testing.T) {
	// Arrange
	handler := handlers.NewGuiHandler()

	// Act
	_, err := handler.LoadState("  \t  ")

	// Assert
	assert.ErrorIs(t, err, common.ErrNoOutputPath)
}

func TestWhenStateFileDoesNotExist_ReturnsNotExistError(t *testing.T) {
	// Arrange
	handler := handlers.NewGuiHandler()
	missingPath := filepath.Join(t.TempDir(), "missing-state.gen.json")

	// Act
	_, err := handler.LoadState(missingPath)

	// Assert
	assert.ErrorIs(t, err, os.ErrNotExist)
}

func TestWhenStateFileContainsInvalidJson_ReturnsError(t *testing.T) {
	// Arrange
	handler := handlers.NewGuiHandler()
	corruptPath := filepath.Join(t.TempDir(), "corrupt-state.gen.json")
	require.NoError(t, os.WriteFile(corruptPath, []byte("this is { not valid json"), 0o644))

	// Act
	_, err := handler.LoadState(corruptPath)

	// Assert
	assert.Error(t, err)
}

func TestWhenStateFileContainsPreviouslySavedState_ReturnsEqualState(t *testing.T) {
	// Arrange
	handler := handlers.NewGuiHandler()
	statePath := filepath.Join(t.TempDir(), "roundtrip-state.gen.json")
	savedState := dtos.NewDefaultEditorStateDto()
	savedState.TemplateName = gofakeit.ProductName()
	_, saveErr := handler.SaveState(dtos.EditorStateSaveDto{
		State:      &savedState,
		OutputPath: statePath,
	})
	require.NoError(t, saveErr)

	// Act
	loadedState, err := handler.LoadState(statePath)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, loadedState)
	assert.Equal(t, savedState, *loadedState)
}
