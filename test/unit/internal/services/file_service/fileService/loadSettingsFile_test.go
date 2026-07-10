package fileService_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/file_service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenSettingsFileIsMissing_ReturnsError(t *testing.T) {
	t.Parallel()
	// Arrange
	missingPath := filepath.Join(t.TempDir(), "missing.gen.json")

	// Act
	_, err := file_service.NewFileService().LoadSettingsFile(missingPath)

	// Assert
	assert.Error(t, err)
}

func TestWhenSettingsFileContainsInvalidJson_ReturnsError(t *testing.T) {
	t.Parallel()
	// Arrange
	settingsPath := filepath.Join(t.TempDir(), "bad.gen.json")
	require.NoError(t, os.WriteFile(settingsPath, []byte("{not json"), 0o644))

	// Act
	_, err := file_service.NewFileService().LoadSettingsFile(settingsPath)

	// Assert
	assert.Error(t, err)
}

func TestWhenSettingsFileContainsValidJson_OverridesPersistedFieldsOnDefaults(t *testing.T) {
	t.Parallel()
	// Arrange
	settingsPath := filepath.Join(t.TempDir(), "ok.gen.json")
	body := `{"templateName":"X","playerCount":4,"mapSize":192}`
	require.NoError(t, os.WriteFile(settingsPath, []byte(body), 0o644))
	expected := dtos.NewDefaultEditorStateDto()
	expected.TemplateName = "X"
	expected.PlayerCount = 4
	expected.MapSize = 192

	// Act
	actual, err := file_service.NewFileService().LoadSettingsFile(settingsPath)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expected, *actual)
}
