package integration_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_errors"
	"github.com/Tariomka/hommoe_custom_templates/internal/composition"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const topologyKey = "topology"

var stateFixtureVersions = []string{ //nolint:gochecknoglobals // Read-only test table.
	"editorState_v0_flat.gen.json",
	"editorState_v1_flat.gen.json",
	"editorState_v2_flat.gen.json",
}

func TestWhenAStateFixtureIsLoadedThroughTheHandler_ItKeepsItsSquareTopology(t *testing.T) {
	t.Parallel()
	for _, name := range stateFixtureVersions {
		t.Run(name+"_KeepsItsSquareTopology", func(t *testing.T) {
			t.Parallel()
			// Arrange
			path := filepath.Join("..", "test_helpers", "testdata", name)

			// Act
			loaded, err := composition.InitializeGuiHandler().LoadState(path, true)

			// Assert
			require.NoError(t, err)
			assert.Equal(t, config.TopologySquare, loaded.State.Topology)
		})
	}
}

func TestWhenAStateOmitsTheTopology_ItLoadsAsRandom(t *testing.T) {
	t.Parallel()
	for _, name := range stateFixtureVersions {
		t.Run(name+"_LoadsAsRandom", func(t *testing.T) {
			t.Parallel()
			// Arrange
			path := writeFixtureWithTopology(t, name, nil)

			// Act
			loaded, err := composition.InitializeGuiHandler().LoadState(path, true)

			// Assert
			require.NoError(t, err)
			assert.Equal(t, config.TopologyRandom, loaded.State.Topology)
		})
	}
}

func TestWhenAStateHasAnEmptyTopology_ItLoadsAsRandom(t *testing.T) {
	t.Parallel()
	for _, name := range stateFixtureVersions {
		t.Run(name+"_LoadsAsRandom", func(t *testing.T) {
			t.Parallel()
			// Arrange
			path := writeFixtureWithTopology(t, name, new(""))

			// Act
			loaded, err := composition.InitializeGuiHandler().LoadState(path, true)

			// Assert
			require.NoError(t, err)
			assert.Equal(t, config.TopologyRandom, loaded.State.Topology)
		})
	}
}

func TestWhenAStateHasAnEmptyTopology_ItWarnsThatRandomIsUsed(t *testing.T) {
	t.Parallel()
	for _, name := range stateFixtureVersions {
		t.Run(name+"_WarnsThatRandomIsUsed", func(t *testing.T) {
			t.Parallel()
			// Arrange
			path := writeFixtureWithTopology(t, name, new(""))

			// Act
			loaded, err := composition.InitializeGuiHandler().LoadState(path, false)

			// Assert
			require.NoError(t, err)
			assert.Contains(t, loaded.Warnings, "topology is empty; using Random")
		})
	}
}

func TestWhenACommittedStateHasARetiredTopology_TheLoadIsRefusedWithItsMessage(t *testing.T) {
	t.Parallel()
	testCases := map[string]string{
		"editorState_v2_retired_Default.gen.json":     retiredMessage("Ring", "Default"),
		"editorState_v2_retired_HubAndSpoke.gen.json": retiredMessage("Hub", "HubAndSpoke"),
		"editorState_v2_retired_Chain.gen.json":       retiredMessage("Chain", "Chain"),
		"editorState_v2_retired_SharedWeb.gen.json":   retiredMessage("Shared Web", "SharedWeb"),
		"editorState_v0_retired_Chain.gen.json":       retiredMessage("Chain", "Chain"),
	}
	for name, expected := range testCases {
		t.Run(name+"_IsRefusedWithItsMessage", func(t *testing.T) {
			t.Parallel()
			// Arrange
			path := filepath.Join("..", "test_helpers", "testdata", name)

			// Act
			_, err := composition.InitializeGuiHandler().LoadState(path, true)

			// Assert
			assert.EqualError(t, err, expected)
		})
	}
}

func TestWhenAV1StateHasARetiredTopology_TheLoadIsRefusedAsUnsupported(t *testing.T) {
	t.Parallel()
	// Arrange
	path := writeFixtureWithTopology(t, "editorState_v1_flat.gen.json", new("HubAndSpoke"))

	// Act
	_, err := composition.InitializeGuiHandler().LoadState(path, true)

	// Assert
	assert.ErrorIs(t, err, common_errors.ErrUnsupportedTopology)
}

func TestWhenAStateHasAnUnknownTopology_TheLoadIsRefusedWithTheUnknownMessage(t *testing.T) {
	t.Parallel()
	// Arrange
	path := writeFixtureWithTopology(t, "editorState_v2_flat.gen.json", new("NotARealTopology"))

	// Act
	_, err := composition.InitializeGuiHandler().LoadState(path, true)

	// Assert
	assert.EqualError(t, err,
		`topology "NotARealTopology" is not a known topology; re-create the template with a supported topology`)
}

func retiredMessage(label, savedID string) string {
	return `topology "` + label + `" (saved as "` + savedID + `") has been retired and is no longer supported; ` +
		"re-create the template with a supported topology"
}

// writeFixtureWithTopology copies a committed fixture into a temp file with its
// topology replaced, or removed when topology is nil.
func writeFixtureWithTopology(t *testing.T, name string, topology *string) string {
	t.Helper()
	keys := parseEditorStateFixtureKeys(t, name)
	delete(keys, topologyKey)
	if topology != nil {
		encoded, err := json.Marshal(*topology)
		require.NoError(t, err)
		keys[topologyKey] = encoded
	}

	content, err := json.Marshal(keys)
	require.NoError(t, err)
	path := filepath.Join(t.TempDir(), name)
	require.NoError(t, os.WriteFile(path, content, 0o644))
	return path
}
