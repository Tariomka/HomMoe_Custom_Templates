//go:build integration_test

package integration_test

import (
	"path/filepath"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/drivers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// retiredStateFixtures pairs a current and a legacy committed fixture with the
// status the editor must show when refusing it.
var retiredStateFixtures = map[string]string{ //nolint:gochecknoglobals // Read-only test table.
	"editorState_v2_retired_Default.gen.json": "Load failed: " + retiredMessage("Ring", "Default") + ".",
	"editorState_v0_retired_Chain.gen.json":   "Load failed: " + retiredMessage("Chain", "Chain") + ".",
}

// documentSnapshot is everything a refused load must leave untouched.
type documentSnapshot struct {
	state    editor_state_model.EditorState
	path     string
	unsaved  bool
	template template_model.Template
}

func TestWhenARetiredStateIsLoaded_TheStatusReportsTheRefusal(t *testing.T) {
	t.Parallel()
	for name, expected := range retiredStateFixtures {
		t.Run(name+"_StatusReportsTheRefusal", func(t *testing.T) {
			t.Parallel()
			// Arrange
			state := newEditedDocument(t)

			// Act
			state.LoadStateFromFile(retiredFixturePath(name))

			// Assert
			message, _ := state.GetStatus()
			assert.Equal(t, expected, message)
		})
	}
}

func TestWhenARetiredStateIsLoaded_TheStatusIsAnError(t *testing.T) {
	t.Parallel()
	for name := range retiredStateFixtures {
		t.Run(name+"_StatusIsAnError", func(t *testing.T) {
			t.Parallel()
			// Arrange
			state := newEditedDocument(t)

			// Act
			state.LoadStateFromFile(retiredFixturePath(name))

			// Assert
			_, isError := state.GetStatus()
			assert.True(t, isError)
		})
	}
}

func TestWhenARetiredStateIsLoaded_TheCurrentDocumentIsUnchanged(t *testing.T) {
	t.Parallel()
	for name := range retiredStateFixtures {
		t.Run(name+"_CurrentDocumentIsUnchanged", func(t *testing.T) {
			t.Parallel()
			// Arrange
			state := newEditedDocument(t)
			expected := snapshotDocument(state)

			// Act
			state.LoadStateFromFile(retiredFixturePath(name))

			// Assert
			assert.Equal(t, expected, snapshotDocument(state))
		})
	}
}

// newEditedDocument returns a session with a generated document that already
// has a path, so every part of it is observable after a refused load.
func newEditedDocument(t *testing.T) *drivers.State {
	t.Helper()
	state := newUIState()
	state.SetCurrentPath(filepath.Join(t.TempDir(), "current.gen.json"))
	state.UpdateState(func(editorState *editor_state_model.EditorState) {
		editorState.TemplateName = "Current Document"
		editorState.Topology = config.TopologySquare
		editorState.PlayerCount = 4
	})
	state.Generate()
	require.NotNil(t, state.GetLastTemplate())
	return state
}

func snapshotDocument(state *drivers.State) documentSnapshot {
	return documentSnapshot{
		state:    state.GetStateData(),
		path:     state.GetCurrentPath(),
		unsaved:  state.IsUnsaved(),
		template: state.GetLastTemplate().Clone(),
	}
}

func retiredFixturePath(name string) string {
	return filepath.Join("..", "test_helpers", "testdata", name)
}
