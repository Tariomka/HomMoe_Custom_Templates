package guiHandler_test

import (
	"path/filepath"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_errors"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenStateToSaveIsNil_ReturnsNothingToSaveError(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	stateSaveDto := editor_state_dto.EditorStateSaveDto{
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
	handler := newProductionGuiHandler()
	state := editor_state_model.NewDefaultEditorStateModel()
	stateSaveDto := editor_state_dto.EditorStateSaveDto{State: toDtoPointer(&state), OutputPath: ""}

	// Act
	_, err := handler.SaveState(stateSaveDto)

	// Assert
	assert.ErrorIs(t, err, common_errors.ErrNoOutputPath)
}

func TestWhenStateOutputPathIsWhitespaceOnly_ReturnsNoOutputPathError(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	state := editor_state_model.NewDefaultEditorStateModel()
	stateSaveDto := editor_state_dto.EditorStateSaveDto{State: toDtoPointer(&state), OutputPath: " \t  "}

	// Act
	_, err := handler.SaveState(stateSaveDto)

	// Assert
	assert.ErrorIs(t, err, common_errors.ErrNoOutputPath)
}

func TestWhenStateAndOutputPathAreValid_ReturnsPathNamedAfterTemplate(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	outputDirectory := t.TempDir()
	state := editor_state_model.NewDefaultEditorStateModel()
	state.TemplateName = "My Template"
	stateSaveDto := editor_state_dto.EditorStateSaveDto{
		State:      toDtoPointer(&state),
		OutputPath: filepath.Join(outputDirectory, "ignored-name.gen.json"),
	}

	// Act
	savedPath, err := handler.SaveState(stateSaveDto)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(outputDirectory, "My Template.gen.json"), savedPath)
}

func TestWhenStateAndOutputPathAreValid_WritesSettingsFile(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	outputPath := filepath.Join(t.TempDir(), "written-state.gen.json")
	state := editor_state_model.NewDefaultEditorStateModel()
	stateSaveDto := editor_state_dto.EditorStateSaveDto{State: toDtoPointer(&state), OutputPath: outputPath}

	// Act
	savedPath, err := handler.SaveState(stateSaveDto)

	// Assert
	require.NoError(t, err)
	assert.FileExists(t, savedPath)
}

func TestWhenStateOutputDirectoryDoesNotExist_CreatesItAndWritesTheFile(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	outputPath := filepath.Join(t.TempDir(), "no-such-directory", "state.gen.json")
	state := editor_state_model.NewDefaultEditorStateModel()
	stateSaveDto := editor_state_dto.EditorStateSaveDto{State: toDtoPointer(&state), OutputPath: outputPath}

	// Act
	savedPath, err := handler.SaveState(stateSaveDto)

	// Assert
	require.NoError(t, err)
	assert.FileExists(t, savedPath)
}
