package settingsFileLoader_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenFileIsMissing_ReturnsError(t *testing.T) {
	// Arrange
	missingPath := filepath.Join(t.TempDir(), "missing.gen.json")

	// Act
	_, err := services.LoadSettingsFile(missingPath)

	// Assert
	assert.Error(t, err)
}

func TestWhenFileContainsInvalidJson_ReturnsError(t *testing.T) {
	// Arrange
	settingsPath := filepath.Join(t.TempDir(), "bad.gen.json")
	require.NoError(t, os.WriteFile(settingsPath, []byte("{not json"), 0o644))

	// Act
	_, err := services.LoadSettingsFile(settingsPath)

	// Assert
	assert.Error(t, err)
}

func TestWhenFileContainsValidJson_OverridesPersistedFieldsOnDefaults(t *testing.T) {
	// Arrange
	settingsPath := filepath.Join(t.TempDir(), "ok.gen.json")
	body := `{"templateName":"X","playerCount":4,"mapSize":192}`
	require.NoError(t, os.WriteFile(settingsPath, []byte(body), 0o644))
	expected := dtos.NewDefaultEditorStateDto()
	expected.TemplateName = "X"
	expected.PlayerCount = 4
	expected.MapSize = 192

	// Act
	actual, err := services.LoadSettingsFile(settingsPath)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expected, *actual)
}
