package guiHandler_test

import (
	"path/filepath"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_errors"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenStateToSaveIsNil_ReturnsNothingToSaveError(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := handlers.NewGuiHandler()
	stateSaveDto := dtos.EditorStateSaveDto{
		State:      nil,
		OutputPath: filepath.Join(t.TempDir(), "state.gen.json"),
	}

	// Act
	_, err := handler.SaveState(stateSaveDto)

	// Assert
	assert.ErrorIs(t, err, common_errors.ErrNothingToSave)
}

func TestWhenStateOutputPathIsEmpty_ReturnsNoOutputPathError(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := handlers.NewGuiHandler()
	state := dtos.NewDefaultEditorStateDto()
	stateSaveDto := dtos.EditorStateSaveDto{State: &state, OutputPath: ""}

	// Act
	_, err := handler.SaveState(stateSaveDto)

	// Assert
	assert.ErrorIs(t, err, common_errors.ErrNoOutputPath)
}

func TestWhenStateOutputPathIsWhitespaceOnly_ReturnsNoOutputPathError(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := handlers.NewGuiHandler()
	state := dtos.NewDefaultEditorStateDto()
	stateSaveDto := dtos.EditorStateSaveDto{State: &state, OutputPath: " \t  "}

	// Act
	_, err := handler.SaveState(stateSaveDto)

	// Assert
	assert.ErrorIs(t, err, common_errors.ErrNoOutputPath)
}

func TestWhenStateAndOutputPathAreValid_ReturnsOutputPath(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := handlers.NewGuiHandler()
	outputPath := filepath.Join(t.TempDir(), "valid-state.gen.json")
	state := dtos.NewDefaultEditorStateDto()
	stateSaveDto := dtos.EditorStateSaveDto{State: &state, OutputPath: outputPath}

	// Act
	savedPath, err := handler.SaveState(stateSaveDto)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, outputPath, savedPath)
}

func TestWhenStateAndOutputPathAreValid_WritesSettingsFile(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := handlers.NewGuiHandler()
	outputPath := filepath.Join(t.TempDir(), "written-state.gen.json")
	state := dtos.NewDefaultEditorStateDto()
	stateSaveDto := dtos.EditorStateSaveDto{State: &state, OutputPath: outputPath}

	// Act
	savedPath, err := handler.SaveState(stateSaveDto)

	// Assert
	require.NoError(t, err)
	assert.FileExists(t, savedPath)
}

func TestWhenStateOutputDirectoryDoesNotExist_ReturnsError(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := handlers.NewGuiHandler()
	outputPath := filepath.Join(t.TempDir(), "no-such-directory", "state.gen.json")
	state := dtos.NewDefaultEditorStateDto()
	stateSaveDto := dtos.EditorStateSaveDto{State: &state, OutputPath: outputPath}

	// Act
	_, err := handler.SaveState(stateSaveDto)

	// Assert
	assert.Error(t, err)
}
