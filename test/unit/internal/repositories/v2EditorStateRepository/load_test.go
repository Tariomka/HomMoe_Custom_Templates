package v2EditorStateRepository_test

import (
	"path/filepath"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state/editor_state_v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenV2StateFileIsMissing_ReturnsError(t *testing.T) {
	t.Parallel()
	// Arrange
	missingPath := filepath.Join(t.TempDir(), "missing.gen.json")
	target := editor_state_v2.EditorState{}

	// Act
	err := newRepository().Load(missingPath, &target)

	// Assert
	assert.Error(t, err)
}

func TestWhenV2StateFileContainsInvalidJson_ReturnsError(t *testing.T) {
	t.Parallel()
	// Arrange
	statePath := writeStateFile(t, "{not json")
	target := editor_state_v2.EditorState{}

	// Act
	err := newRepository().Load(statePath, &target)

	// Assert
	assert.Error(t, err)
}

// The flat layout and the text ban list are what the current entity can no
// longer read, so they are what this repository has to decode.
func TestWhenV2StateFileCarriesFlatKeys_TheyReachTheTarget(t *testing.T) {
	t.Parallel()
	// Arrange
	statePath := writeStateFile(t, `{"schemaVersion":2,"playerCount":4,"bannedItems":"a\nb"}`)
	target := editor_state_v2.EditorState{}

	// Act
	err := newRepository().Load(statePath, &target)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, []any{4, "a\nb"}, []any{target.PlayerCount, target.BannedItems})
}

func TestWhenV2StateFileOmitsAKey_TheValueSeededByTheCallerSurvives(t *testing.T) {
	t.Parallel()
	// Arrange
	statePath := writeStateFile(t, `{"schemaVersion":2,"playerCount":4}`)
	target := editor_state_v2.EditorState{
		TemplateName: "Seeded"}

	// Act
	err := newRepository().Load(statePath, &target)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "Seeded", target.TemplateName)
}
