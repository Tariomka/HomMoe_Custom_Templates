package integration_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/composition"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// A migrated document saved as v3 and loaded again must be the same document:
// every setting and every manual zone and connection, nil and empty included.
func TestWhenACuratedV2DocumentIsSavedAsV3AndReloaded_TheStateIsUnchanged(t *testing.T) {
	t.Parallel()
	for _, fixture := range curatedManualEditFixtures() {
		t.Run(fixture+"_ReloadsUnchanged", func(t *testing.T) {
			t.Parallel()
			// Arrange
			handler := composition.InitializeGuiHandler()
			migrated, err := handler.LoadState(filepath.Join("..", "test_helpers", "testdata", fixture), false)
			require.NoError(t, err)
			require.True(t, migrated.State.HasManualEdits())
			savedPath := saveStateTo(t, migrated.State, t.TempDir())

			// Act
			reloaded, err := handler.LoadState(savedPath, false)

			// Assert
			require.NoError(t, err)
			assert.Equal(t, migrated.State, reloaded.State)
		})
	}
}

// Saving a reloaded v3 document again must write the same bytes, otherwise
// every save would churn the file.
func TestWhenACuratedDocumentIsSavedTwice_TheBytesAreStable(t *testing.T) {
	t.Parallel()
	for _, fixture := range curatedManualEditFixtures() {
		t.Run(fixture+"_IsByteStable", func(t *testing.T) {
			t.Parallel()
			// Arrange
			handler := composition.InitializeGuiHandler()
			migrated, err := handler.LoadState(filepath.Join("..", "test_helpers", "testdata", fixture), false)
			require.NoError(t, err)
			firstPath := saveStateTo(t, migrated.State, t.TempDir())
			reloaded, err := handler.LoadState(firstPath, false)
			require.NoError(t, err)

			// Act
			secondPath := saveStateTo(t, reloaded.State, t.TempDir())

			// Assert
			assert.Equal(t, readFile(t, firstPath), readFile(t, secondPath))
		})
	}
}

// curatedManualEditFixtures are real v2 documents built by the Batch L probe
// through the driver: one per kind of manual edit on a deterministic layout,
// plus Random-topology and 6-player hub layouts.
func curatedManualEditFixtures() []string {
	return []string{
		"editorState_v2_square2p_move.gen.json",
		"editorState_v2_square2p_add_zone.gen.json",
		"editorState_v2_square2p_delete_zone.gen.json",
		"editorState_v2_square2p_add_connection.gen.json",
		"editorState_v2_square2p_delete_connection.gen.json",
		"editorState_v2_square2p_edit_connection.gen.json",
		"editorState_v2_square2p_quality.gen.json",
		"editorState_v2_square2p_field_override.gen.json",
		"editorState_v2_square2p_nil_empty.gen.json",
		"editorState_v2_random8p_move.gen.json",
		"editorState_v2_random8p_add_zone.gen.json",
		"editorState_v2_random8p_nil_empty.gen.json",
		"editorState_v2_geohub6p_move.gen.json",
		"editorState_v2_geohub6p_delete_zone.gen.json",
	}
}

func saveStateTo(t *testing.T, state editor_state_model.EditorState, directory string) string {
	t.Helper()

	savedPath, err := composition.InitializeGuiHandler().SaveState(editor_state_dto.EditorStateSaveDto{
		State:      &editor_state_dto.EditorStateDto{EditorState: state},
		OutputPath: filepath.Join(directory, "state.gen.json"),
	})
	require.NoError(t, err)

	return savedPath
}

func readFile(t *testing.T, path string) string {
	t.Helper()

	data, err := os.ReadFile(path)
	require.NoError(t, err)

	return string(data)
}
