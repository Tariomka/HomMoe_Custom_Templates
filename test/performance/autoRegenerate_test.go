package performance_test

import (
	"testing"
	"time"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/drivers"
	"github.com/Tariomka/hommoe_custom_templates/internal/composition"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/stretchr/testify/require"
)

// BenchmarkAutoRegenerate_IdleFrame measures the regeneration check the editor
// window runs on every frame, once the template is generated and nothing has
// changed since.
func BenchmarkAutoRegenerate_IdleFrame(b *testing.B) {
	withEdits := loadLargeManualEditState(b)
	withoutEdits := withEdits.Clone()
	withoutEdits.ManualEditSettings = editor_state_model.ManualEditSettings{}

	benchmarkCases := []struct {
		name  string
		state editor_state_model.EditorState
	}{
		{name: "WithManualEdits", state: withEdits},
		{name: "WithoutManualEdits", state: withoutEdits},
	}

	for _, benchmarkCase := range benchmarkCases {
		b.Run(benchmarkCase.name, func(b *testing.B) {
			state := newGeneratedState(b, benchmarkCase.state)
			now := time.Now()

			b.ReportAllocs()
			for b.Loop() {
				state.AutoRegenerate(now)
			}

			currentState := state.GetStateData()
			require.Equal(b, benchmarkCase.state.HasManualEdits(), currentState.HasManualEdits())
		})
	}
}

// newGeneratedState builds a window-free driver on the production object graph,
// puts the given state in place and runs the first generation, so the next
// AutoRegenerate call is an idle frame.
func newGeneratedState(tb testing.TB, editorState editor_state_model.EditorState) *drivers.State {
	tb.Helper()
	state := drivers.NewUIState(
		composition.InitializeGuiHandler(),
		composition.InitializeFileSystemHandler(),
		composition.InitializeRegenerationHandler(),
		false)
	state.UpdateState(func(current *editor_state_model.EditorState) { *current = editorState.Clone() })
	state.AutoRegenerate(time.Now())
	require.NotNil(tb, state.GetLastTemplate())
	return state
}
