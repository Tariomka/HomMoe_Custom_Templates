package fileService_test

import (
	"math"
	"os"
	"path/filepath"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/file_service"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenStateIsSaved_WritesIndentedJson(t *testing.T) {
	// Arrange
	settingsPath := filepath.Join(t.TempDir(), "out.gen.json")
	state := dtos.NewDefaultEditorStateDto()
	state.TemplateName = gofakeit.ProductName()
	require.NoError(t, file_service.NewFileService().SaveSettings(settingsPath, &state))

	// Act
	data, err := os.ReadFile(settingsPath)

	// Assert
	require.NoError(t, err)
	assert.Contains(t, string(data), "\n  ")
}

func TestWhenSavedStateIsLoaded_RoundTripsState(t *testing.T) {
	// Arrange
	settingsPath := filepath.Join(t.TempDir(), "roundtrip.gen.json")
	service := file_service.NewFileService()
	state := dtos.NewDefaultEditorStateDto()
	state.TemplateName = gofakeit.ProductName()
	state.PlayerCount = gofakeit.Number(2, 8)
	require.NoError(t, service.SaveSettings(settingsPath, &state))

	// Act
	loaded, err := service.LoadSettingsFile(settingsPath)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, state, *loaded)
}

func TestWhenSettingsDirectoryDoesNotExist_ReturnsError(t *testing.T) {
	// Arrange
	missingDirPath := filepath.Join(t.TempDir(), "missing_dir", "x.gen.json")
	state := dtos.NewDefaultEditorStateDto()

	// Act
	err := file_service.NewFileService().SaveSettings(missingDirPath, &state)

	// Assert
	assert.Error(t, err)
}

func TestWhenStateContainsNaNValue_ReturnsError(t *testing.T) {
	// Arrange
	settingsPath := filepath.Join(t.TempDir(), "nan.gen.json")
	state := dtos.NewDefaultEditorStateDto()
	state.PlayerZoneSize = math.NaN()

	// Act
	err := file_service.NewFileService().SaveSettings(settingsPath, &state)

	// Assert
	assert.Error(t, err)
}
