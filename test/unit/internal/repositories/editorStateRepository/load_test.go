package editorStateRepository_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenStateFileIsMissing_ReturnsError(t *testing.T) {
	t.Parallel()
	// Arrange
	missingPath := filepath.Join(t.TempDir(), "missing.gen.json")
	target := editor_state.EditorState{}

	// Act
	err := newRepository().Load(missingPath, &target)

	// Assert
	assert.Error(t, err)
}

func TestWhenStateFileContainsInvalidJson_ReturnsError(t *testing.T) {
	t.Parallel()
	// Arrange
	statePath := filepath.Join(t.TempDir(), "bad.gen.json")
	require.NoError(t, os.WriteFile(statePath, []byte("{not json"), 0o644))
	target := editor_state.EditorState{}

	// Act
	err := newRepository().Load(statePath, &target)

	// Assert
	assert.Error(t, err)
}

func TestWhenStateFileContainsValidJson_TheDecodedKeysReachTheTarget(t *testing.T) {
	t.Parallel()
	// Arrange
	statePath := writeStateFile(t, `{"templateName":"X","playerCount":4,"mapSize":192}`)
	expected := editor_state.EditorState{TemplateName: "X", PlayerCount: 4, MapSize: 192}
	target := editor_state.EditorState{}

	// Act
	err := newRepository().Load(statePath, &target)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expected, target)
}

// The repository decodes *into* the value it is handed and never seeds one, so
// whatever the caller put there survives the keys the file omits. FileService
// relies on this to hand the load the default entity.
func TestWhenStateFileOmitsAKey_TheValueSeededByTheCallerSurvives(t *testing.T) {
	t.Parallel()
	// Arrange
	statePath := writeStateFile(t, `{"playerCount":4}`)
	target := editor_state.EditorState{TemplateName: "Seeded"}

	// Act
	err := newRepository().Load(statePath, &target)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "Seeded", target.TemplateName)
}

func TestWhenStateFileOmitsTheContentRows_TheRowsSeededByTheCallerSurvive(t *testing.T) {
	t.Parallel()
	// Arrange
	statePath := writeStateFile(t, `{"templateName":"X"}`)
	seededRows := []editor_state.ZoneContentRow{{Sid: "seeded", Count: 1}}
	target := editor_state.EditorState{PlayerZoneContentRows: seededRows}

	// Act
	err := newRepository().Load(statePath, &target)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, seededRows, target.PlayerZoneContentRows)
}

// Written empty and written nil are the same absent key on disk, so neither can
// clear what the caller seeded. That is a property of the `omitempty` tags, not
// a defect of the loader - the nil-versus-empty distinction the change
// detection draws never survived a save and is not meant to.
func TestWhenStateFileCarriesEmptyContentRows_TheRowsSeededByTheCallerSurviveToo(t *testing.T) {
	t.Parallel()
	// Arrange
	repository := newRepository()
	saved := editor_state.EditorState{PlayerZoneContentRows: []editor_state.ZoneContentRow{}}
	statePath, err := repository.Save(t.TempDir(), "State", saved)
	require.NoError(t, err)
	seededRows := []editor_state.ZoneContentRow{{Sid: "seeded", Count: 1}}
	target := editor_state.EditorState{PlayerZoneContentRows: seededRows}

	// Act
	loadErr := repository.Load(statePath, &target)

	// Assert
	require.NoError(t, loadErr)
	assert.Equal(t, seededRows, target.PlayerZoneContentRows)
}

func writeStateFile(t *testing.T, body string) string {
	t.Helper()

	statePath := filepath.Join(t.TempDir(), "state.gen.json")
	require.NoError(t, os.WriteFile(statePath, []byte(body), 0o644))

	return statePath
}
