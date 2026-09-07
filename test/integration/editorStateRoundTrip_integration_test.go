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
	first, err := handler.LoadState(savedPath, false)
	require.NoError(t, err)
	first.State.TemplateName = "Mutated After Load"

	// Act
	second, err := handler.LoadState(savedPath, false)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, test_helpers.NewAllFieldsEditorStateModel().TemplateName, second.State.TemplateName)
}

func TestWhenALoadedStateSliceIsMutated_ReloadingIsUnaffected(t *testing.T) {
	t.Parallel()
	// Arrange
	handler, savedPath := saveRoundTripState(t)
	first, err := handler.LoadState(savedPath, false)
	require.NoError(t, err)
	require.NotEmpty(t, first.State.PlayerZoneContentRows)
	first.State.PlayerZoneContentRows[0].Sid = "mutated-after-load"

	// Act
	second, err := handler.LoadState(savedPath, false)

	// Assert
	require.NoError(t, err)
	assert.Equal(
		t,
		test_helpers.NewAllFieldsEditorStateModel().PlayerZoneContentRows,
		second.State.PlayerZoneContentRows)
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
	loaded, err := handler.LoadState(savedPath, false)

	// Assert
	require.NoError(t, err)
	assert.Equal(
		t,
		test_helpers.NewAllFieldsEditorStateModel().PlayerZoneContentRows,
		loaded.State.PlayerZoneContentRows)
}

func TestWhenAStateIsRoundTripped_TheReloadedModelEqualsTheSavedOne(t *testing.T) {
	t.Parallel()
	// Arrange
	handler, savedPath := saveRoundTripState(t)

	// Act
	loaded, err := handler.LoadState(savedPath, false)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, test_helpers.NewAllFieldsEditorStateModel(), loaded.State)
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
