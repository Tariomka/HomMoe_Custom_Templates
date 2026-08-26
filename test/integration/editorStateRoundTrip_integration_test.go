package integration_test

import (
	"path/filepath"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/composition"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers/handler_interfaces"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenALoadedStateScalarIsMutated_ReloadingIsUnaffected(t *testing.T) {
	t.Parallel()
	// Arrange
	handler, savedPath := saveRoundTripState(t)
	first, _, err := handler.LoadState(savedPath, false)
	require.NoError(t, err)
	first.TemplateName = "Mutated After Load"

	// Act
	second, _, err := handler.LoadState(savedPath, false)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, test_helpers.NewAllFieldsEditorStateModel().TemplateName, second.TemplateName)
}

func TestWhenALoadedStateSliceIsMutated_ReloadingIsUnaffected(t *testing.T) {
	t.Parallel()
	// Arrange
	handler, savedPath := saveRoundTripState(t)
	first, _, err := handler.LoadState(savedPath, false)
	require.NoError(t, err)
	require.NotEmpty(t, first.PlayerZoneContentRows)
	first.PlayerZoneContentRows[0].Sid = "mutated-after-load"

	// Act
	second, _, err := handler.LoadState(savedPath, false)

	// Assert
	require.NoError(t, err)
	assert.Equal(
		t,
		test_helpers.NewAllFieldsEditorStateModel().PlayerZoneContentRows,
		second.PlayerZoneContentRows)
}

func TestWhenAStateIsSavedThroughTheHandler_TheSourceModelIsNotAliasedByTheFile(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := composition.InitializeGuiHandler()
	state := test_helpers.NewAllFieldsEditorStateModel()
	savedPath, err := handler.SaveState(editor_state_dto.EditorStateSaveDto{
		State:      &editor_state_dto.EditorStateDto{EditorState: state},
		OutputPath: filepath.Join(t.TempDir(), "ignored.gen.json"),
	})
	require.NoError(t, err)
	state.PlayerZoneContentRows[0].Sid = "mutated-after-save"

	// Act
	loaded, _, err := handler.LoadState(savedPath, false)

	// Assert
	require.NoError(t, err)
	assert.Equal(
		t,
		test_helpers.NewAllFieldsEditorStateModel().PlayerZoneContentRows,
		loaded.PlayerZoneContentRows)
}

func TestWhenAStateIsRoundTripped_TheReloadedModelEqualsTheSavedOne(t *testing.T) {
	t.Parallel()
	// Arrange
	handler, savedPath := saveRoundTripState(t)

	// Act
	loaded, _, err := handler.LoadState(savedPath, false)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, editor_state_dto.EditorStateDto{EditorState: test_helpers.NewAllFieldsEditorStateModel()}, *loaded)
}

// saveRoundTripState writes the all-fields state through the real handler graph
// and returns the handler together with the path it wrote.
func saveRoundTripState(t *testing.T) (handler_interfaces.IGuiHandler, string) {
	t.Helper()
	handler := composition.InitializeGuiHandler()
	state := test_helpers.NewAllFieldsEditorStateModel()
	savedPath, err := handler.SaveState(editor_state_dto.EditorStateSaveDto{
		State:      &editor_state_dto.EditorStateDto{EditorState: state},
		OutputPath: filepath.Join(t.TempDir(), "ignored.gen.json"),
	})
	require.NoError(t, err)

	return handler, savedPath
}
