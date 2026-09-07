package legacyEditorStateRepository_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state/editor_state_v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenLegacyStateFileIsMissing_ReturnsError(t *testing.T) {
	t.Parallel()
	// Arrange
	missingPath := filepath.Join(t.TempDir(), "missing.gen.json")
	target := editor_state_v1.EditorState{}

	// Act
	err := newRepository().Load(missingPath, &target)

	// Assert
	assert.Error(t, err)
}

func TestWhenLegacyStateFileContainsInvalidJson_ReturnsError(t *testing.T) {
	t.Parallel()
	// Arrange
	statePath := writeStateFile(t, "{not json")
	target := editor_state_v1.EditorState{}

	// Act
	err := newRepository().Load(statePath, &target)

	// Assert
	assert.Error(t, err)
}

// The array-shaped manualPosition is the whole reason this repository exists:
// the current entity cannot decode it at all.
func TestWhenLegacyStateFileCarriesAnArrayPosition_ItReachesTheTarget(t *testing.T) {
	t.Parallel()
	// Arrange
	statePath := writeStateFile(
		t,
		`{"schemaVersion":1,"manualZones":[{"zone":{"name":"A"},"manualPosition":[0.25,0.75]}]}`)
	target := editor_state_v1.EditorState{}

	// Act
	err := newRepository().Load(statePath, &target)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, [2]float64{0.25, 0.75}, *target.ManualZones[0].ManualPosition)
}

func TestWhenLegacyStateFileOmitsAKey_TheValueSeededByTheCallerSurvives(t *testing.T) {
	t.Parallel()
	// Arrange
	statePath := writeStateFile(t, `{"schemaVersion":1,"playerCount":4}`)
	target := editor_state_v1.EditorState{TemplateName: "Seeded"}

	// Act
	err := newRepository().Load(statePath, &target)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "Seeded", target.TemplateName)
}

func writeStateFile(t *testing.T, body string) string {
	t.Helper()

	statePath := filepath.Join(t.TempDir(), "state.gen.json")
	require.NoError(t, os.WriteFile(statePath, []byte(body), 0o644))

	return statePath
}
