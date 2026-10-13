package performance_test

import (
	"path/filepath"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/composition"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/stretchr/testify/require"
)

const largeManualEditFixture = "editorState_v2_geohub6p_move.gen.json"

// BenchmarkEditorStateFile_SaveLoad measures saving and reloading a 6-player
// layout that carries manual zone edits, through the production state handler.
func BenchmarkEditorStateFile_SaveLoad(b *testing.B) {
	handler := composition.InitializeGuiHandler()
	state := editor_state_dto.EditorStateDto{EditorState: loadLargeManualEditState(b)}
	outputPath := filepath.Join(b.TempDir(), "bench.gen.json")

	b.ReportAllocs()
	for b.Loop() {
		savedPath, err := handler.SaveState(editor_state_dto.EditorStateSaveDto{State: &state, OutputPath: outputPath})
		require.NoError(b, err)
		_, err = handler.LoadState(savedPath, false)
		require.NoError(b, err)
	}
}

// loadLargeManualEditState loads the committed 6-player fixture, whose manual
// edits snapshot every zone and connection of the layout.
func loadLargeManualEditState(tb testing.TB) editor_state_model.EditorState {
	tb.Helper()
	path := filepath.Join("..", "test_helpers", "testdata", largeManualEditFixture)
	loaded, err := composition.InitializeGuiHandler().LoadState(path, true)
	require.NoError(tb, err)
	require.True(tb, loaded.State.HasManualEdits())
	return loaded.State
}
