package editorStateRepository_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenStateFileIsMissing_ReturnsError(t *testing.T) {
	t.Parallel()
	// Arrange
	missingPath := filepath.Join(t.TempDir(), "missing.gen.json")

	// Act
	_, err := newRepository().Load(missingPath)

	// Assert
	assert.Error(t, err)
}

func TestWhenStateFileContainsInvalidJson_ReturnsError(t *testing.T) {
	t.Parallel()
	// Arrange
	statePath := filepath.Join(t.TempDir(), "bad.gen.json")
	require.NoError(t, os.WriteFile(statePath, []byte("{not json"), 0o644))

	// Act
	_, err := newRepository().Load(statePath)

	// Assert
	assert.Error(t, err)
}

func TestWhenStateFileContainsValidJson_OverridesPersistedFieldsOnDefaults(t *testing.T) {
	t.Parallel()
	// Arrange
	statePath := filepath.Join(t.TempDir(), "ok.gen.json")
	body := `{"templateName":"X","playerCount":4,"mapSize":192}`
	require.NoError(t, os.WriteFile(statePath, []byte(body), 0o644))
	expected := editor_state_model.NewDefaultEditorStateModel()
	expected.TemplateName = "X"
	expected.PlayerCount = 4
	expected.MapSize = 192

	// Act
	actual, err := newRepository().Load(statePath)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestWhenStateFileOmitsTheContentRows_TheSeededDefaultRowsSurvive(t *testing.T) {
	t.Parallel()
	// Arrange
	statePath := filepath.Join(t.TempDir(), "noRows.gen.json")
	require.NoError(t, os.WriteFile(statePath, []byte(`{"templateName":"X"}`), 0o644))

	// Act
	actual, err := newRepository().Load(statePath)

	// Assert
	require.NoError(t, err)
	assert.Equal(
		t,
		editor_state_model.NewDefaultEditorStateModel().PlayerZoneContentRows,
		actual.PlayerZoneContentRows)
}

// Written empty and written nil are the same absent key on disk, so both come
// back as the seeded default. That is a property of the `omitempty` tags, not a
// defect of the loader - the nil-versus-empty distinction the change detection
// draws never survived a save and is not meant to.
func TestWhenStateFileCarriesEmptyContentRows_TheSeededDefaultRowsSurviveToo(t *testing.T) {
	t.Parallel()
	// Arrange
	repository := newRepository()
	saved := editor_state_model.NewDefaultEditorStateModel()
	saved.PlayerZoneContentRows = []models.ZoneContentRow{}
	statePath, err := repository.Save(t.TempDir(), "State", saved)
	require.NoError(t, err)

	// Act
	actual, loadErr := repository.Load(statePath)

	// Assert
	require.NoError(t, loadErr)
	assert.Equal(
		t,
		editor_state_model.NewDefaultEditorStateModel().PlayerZoneContentRows,
		actual.PlayerZoneContentRows)
}
