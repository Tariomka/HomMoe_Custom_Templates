package editorStateRepository_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/repositories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenStateFileIsMissing_ReturnsError(t *testing.T) {
	t.Parallel()
	// Arrange
	missingPath := filepath.Join(t.TempDir(), "missing.gen.json")

	// Act
	_, err := repositories.NewEditorStateRepository().Load(missingPath)

	// Assert
	assert.Error(t, err)
}

func TestWhenStateFileContainsInvalidJson_ReturnsError(t *testing.T) {
	t.Parallel()
	// Arrange
	statePath := filepath.Join(t.TempDir(), "bad.gen.json")
	require.NoError(t, os.WriteFile(statePath, []byte("{not json"), 0o644))

	// Act
	_, err := repositories.NewEditorStateRepository().Load(statePath)

	// Assert
	assert.Error(t, err)
}

func TestWhenStateFileContainsValidJson_OverridesPersistedFieldsOnDefaults(t *testing.T) {
	t.Parallel()
	// Arrange
	statePath := filepath.Join(t.TempDir(), "ok.gen.json")
	body := `{"templateName":"X","playerCount":4,"mapSize":192}`
	require.NoError(t, os.WriteFile(statePath, []byte(body), 0o644))
	expected := editor_state_dto.NewDefaultEditorStateDto()
	expected.TemplateName = "X"
	expected.PlayerCount = 4
	expected.MapSize = 192

	// Act
	actual, err := repositories.NewEditorStateRepository().Load(statePath)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expected, actual)
}
