package guiHandler_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_errors"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenStateFilePathIsEmpty_ReturnsNoOutputPathError(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()

	// Act
	_, _, err := handler.LoadState("", true)

	// Assert
	assert.ErrorIs(t, err, common_errors.ErrNoOutputPath)
}

func TestWhenStateFilePathIsWhitespaceOnly_ReturnsNoOutputPathError(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()

	// Act
	_, _, err := handler.LoadState("  \t  ", true)

	// Assert
	assert.ErrorIs(t, err, common_errors.ErrNoOutputPath)
}

func TestWhenStateFileDoesNotExist_ReturnsNotExistError(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	missingPath := filepath.Join(t.TempDir(), "missing-state.gen.json")

	// Act
	_, _, err := handler.LoadState(missingPath, true)

	// Assert
	assert.ErrorIs(t, err, os.ErrNotExist)
}

func TestWhenStateFileContainsInvalidJson_ReturnsError(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	corruptPath := filepath.Join(t.TempDir(), "corrupt-state.gen.json")
	require.NoError(t, os.WriteFile(corruptPath, []byte("this is { not valid json"), 0o644))

	// Act
	_, _, err := handler.LoadState(corruptPath, true)

	// Assert
	assert.Error(t, err)
}

func TestWhenStateFileContainsPreviouslySavedState_ReturnsEqualState(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	statePath := filepath.Join(t.TempDir(), "roundtrip-state.gen.json")
	savedState := editor_state_model.NewDefaultEditorStateModel()
	savedState.TemplateName = gofakeit.ProductName()
	savedPath, saveErr := handler.SaveState(editor_state_dto.EditorStateSaveDto{
		State:      &savedState,
		OutputPath: statePath,
	})
	require.NoError(t, saveErr)

	// Act
	loadedState, _, err := handler.LoadState(savedPath, true)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, loadedState)
	assert.Equal(t, savedState, *loadedState)
}

func TestWhenStateFileIsValid_ReturnsNoWarnings(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	statePath := filepath.Join(t.TempDir(), "valid-state.gen.json")
	savedState := editor_state_model.NewDefaultEditorStateModel()
	savedPath, saveErr := handler.SaveState(editor_state_dto.EditorStateSaveDto{
		State:      &savedState,
		OutputPath: statePath,
	})
	require.NoError(t, saveErr)

	// Act
	_, warnings, err := handler.LoadState(savedPath, true)

	// Assert
	require.NoError(t, err)
	assert.Empty(t, warnings)
}

func TestWhenStateFileHasOutOfRangeValues_ReturnsWarnings(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	statePath := filepath.Join(t.TempDir(), "invalid-state.gen.json")
	require.NoError(t, os.WriteFile(statePath, []byte(`{"playerCount": 50}`), 0o644))

	// Act
	_, warnings, err := handler.LoadState(statePath, true)

	// Assert
	require.NoError(t, err)
	assert.NotEmpty(t, warnings)
}

func TestWhenFixIssuesIsTrue_ReturnsStateWithIssuesFixed(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	statePath := filepath.Join(t.TempDir(), "invalid-state.gen.json")
	require.NoError(t, os.WriteFile(statePath, []byte(`{"playerCount": 50}`), 0o644))

	// Act
	loadedState, _, err := handler.LoadState(statePath, true)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, 8, loadedState.PlayerCount)
}

func TestWhenFixIssuesIsFalse_ReturnsStateWithIssuesUnfixed(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	statePath := filepath.Join(t.TempDir(), "invalid-state.gen.json")
	require.NoError(t, os.WriteFile(statePath, []byte(`{"playerCount": 50}`), 0o644))

	// Act
	loadedState, _, err := handler.LoadState(statePath, false)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, 50, loadedState.PlayerCount)
}

func TestWhenFixIssuesIsFalse_StillReturnsWarnings(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	statePath := filepath.Join(t.TempDir(), "invalid-state.gen.json")
	require.NoError(t, os.WriteFile(statePath, []byte(`{"playerCount": 50}`), 0o644))

	// Act
	_, warnings, err := handler.LoadState(statePath, false)

	// Assert
	require.NoError(t, err)
	assert.NotEmpty(t, warnings)
}
