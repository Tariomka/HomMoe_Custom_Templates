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
	_, err := handler.LoadState("", true)

	// Assert
	assert.ErrorIs(t, err, common_errors.ErrNoOutputPath)
}

func TestWhenStateFilePathIsWhitespaceOnly_ReturnsNoOutputPathError(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()

	// Act
	_, err := handler.LoadState("  \t  ", true)

	// Assert
	assert.ErrorIs(t, err, common_errors.ErrNoOutputPath)
}

func TestWhenStateFileDoesNotExist_ReturnsNotExistError(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	missingPath := filepath.Join(t.TempDir(), "missing-state.gen.json")

	// Act
	_, err := handler.LoadState(missingPath, true)

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
	_, err := handler.LoadState(corruptPath, true)

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
		State:      toDtoPointer(&savedState),
		OutputPath: statePath,
	})
	require.NoError(t, saveErr)

	// Act
	loaded, err := handler.LoadState(savedPath, true)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, loaded)
	assert.Equal(t, savedState, loaded.State)
}

func TestWhenStateFileIsValid_ReturnsNoWarnings(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	statePath := filepath.Join(t.TempDir(), "valid-state.gen.json")
	savedState := editor_state_model.NewDefaultEditorStateModel()
	savedPath, saveErr := handler.SaveState(editor_state_dto.EditorStateSaveDto{
		State:      toDtoPointer(&savedState),
		OutputPath: statePath,
	})
	require.NoError(t, saveErr)

	// Act
	loaded, err := handler.LoadState(savedPath, true)

	// Assert
	require.NoError(t, err)
	assert.Empty(t, loaded.Warnings)
}

func TestWhenStateFileHasOutOfRangeValues_ReturnsWarnings(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	statePath := filepath.Join(t.TempDir(), "invalid-state.gen.json")
	require.NoError(t, os.WriteFile(statePath, []byte(`{"playerCount": 50}`), 0o644))

	// Act
	loaded, err := handler.LoadState(statePath, true)

	// Assert
	require.NoError(t, err)
	assert.NotEmpty(t, loaded.Warnings)
}

func TestWhenFixIssuesIsTrue_ReturnsStateWithIssuesFixed(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	statePath := filepath.Join(t.TempDir(), "invalid-state.gen.json")
	require.NoError(t, os.WriteFile(statePath, []byte(`{"playerCount": 50}`), 0o644))

	// Act
	loaded, err := handler.LoadState(statePath, true)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, 8, loaded.State.PlayerCount)
}

func TestWhenFixIssuesIsFalse_ReturnsStateWithIssuesUnfixed(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	statePath := filepath.Join(t.TempDir(), "invalid-state.gen.json")
	require.NoError(t, os.WriteFile(statePath, []byte(`{"playerCount": 50}`), 0o644))

	// Act
	loaded, err := handler.LoadState(statePath, false)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, 50, loaded.State.PlayerCount)
}

func TestWhenFixIssuesIsFalse_StillReturnsWarnings(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	statePath := filepath.Join(t.TempDir(), "invalid-state.gen.json")
	require.NoError(t, os.WriteFile(statePath, []byte(`{"playerCount": 50}`), 0o644))

	// Act
	loaded, err := handler.LoadState(statePath, false)

	// Assert
	require.NoError(t, err)
	assert.NotEmpty(t, loaded.Warnings)
}

func TestWhenCommittedStateHasARetiredTopology_ReturnsUnsupportedTopologyError(t *testing.T) {
	t.Parallel()
	for _, name := range []string{
		"editorState_v2_retired_Default.gen.json",
		"editorState_v2_retired_HubAndSpoke.gen.json",
		"editorState_v2_retired_Chain.gen.json",
		"editorState_v2_retired_SharedWeb.gen.json",
		"editorState_v0_retired_Chain.gen.json",
	} {
		t.Run(name+"_ReturnsUnsupportedTopologyError", func(t *testing.T) {
			t.Parallel()
			// Arrange
			handler := newProductionGuiHandler()
			statePath := filepath.Join("..", "..", "..", "..", "test_helpers", "testdata", name)

			// Act
			_, err := handler.LoadState(statePath, false)

			// Assert
			assert.ErrorIs(t, err, common_errors.ErrUnsupportedTopology)
		})
	}
}

func TestWhenStateFileHasAnUnknownTopology_ReturnsUnsupportedTopologyError(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	statePath := filepath.Join(t.TempDir(), "unknown-topology.gen.json")
	require.NoError(t, os.WriteFile(statePath, []byte(`{"schemaVersion": 2, "topology": "NotARealTopology"}`), 0o644))

	// Act
	_, err := handler.LoadState(statePath, true)

	// Assert
	assert.ErrorIs(t, err, common_errors.ErrUnsupportedTopology)
}
