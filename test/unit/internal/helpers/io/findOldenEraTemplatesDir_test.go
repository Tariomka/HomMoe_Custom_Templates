package io_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_errors"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenUserTemplatesFolderExists_ReturnsTemplatesPath(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("user-directory glob detection runs only on Windows")
	}
	// Arrange
	homeDir := t.TempDir()
	t.Setenv("USERPROFILE", homeDir)
	expectedPath := filepath.Join(
		homeDir, "AppData", "LocalLow", "Unfrozen", "HeroesOldenEra",
		"users", gofakeit.LetterN(8), "my_map_templates")
	require.NoError(t, os.MkdirAll(expectedPath, 0o750))

	// Act
	templatesPath, err := helpers.FindOldenEraTemplatesDir(false)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedPath, templatesPath)
}

func TestWhenUserTemplatesFolderIsMissing_ReturnsTemplatesDirNotFoundError(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("user-directory glob detection runs only on Windows")
	}
	// Arrange
	t.Setenv("USERPROFILE", t.TempDir())

	// Act
	_, err := helpers.FindOldenEraTemplatesDir(false)

	// Assert
	assert.ErrorIs(t, err, common_errors.ErrTemplatesDirNotFound)
}

func TestWhenSteamLibraryContainsGameAndProtonPrefixExists_ReturnsTemplatesPath(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Steam VDF detection for user templates runs only on non-Windows")
	}
	// Arrange
	homeDir := t.TempDir()
	t.Setenv("HOME", homeDir)
	libraryDir := writeSteamLibraryVDF(t, homeDir, "3105440")
	expectedPath := filepath.Join(
		libraryDir, "steamapps", "compatdata", "3105440", "pfx", "drive_c", "users", "steamuser",
		"AppData", "LocalLow", "Unfrozen", "HeroesOldenEra",
		"users", gofakeit.LetterN(8), "my_map_templates")
	require.NoError(t, os.MkdirAll(expectedPath, 0o750))

	// Act
	templatesPath, err := helpers.FindOldenEraTemplatesDir(false)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedPath, templatesPath)
}

func TestWhenSteamLibraryContainsGameButProtonPrefixIsMissing_ReturnsTemplatesDirNotFoundError(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Steam VDF detection for user templates runs only on non-Windows")
	}
	// Arrange
	homeDir := t.TempDir()
	t.Setenv("HOME", homeDir)
	writeSteamLibraryVDF(t, homeDir, "3105440")

	// Act
	_, err := helpers.FindOldenEraTemplatesDir(false)

	// Assert
	assert.ErrorIs(t, err, common_errors.ErrTemplatesDirNotFound)
}

func TestWhenSteamLibraryDoesNotContainGame_ReturnsGameInVDFNotFoundError(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Steam VDF detection for user templates runs only on non-Windows")
	}
	// Arrange
	homeDir := t.TempDir()
	t.Setenv("HOME", homeDir)
	writeSteamLibraryVDF(t, homeDir, gofakeit.DigitN(6))

	// Act
	_, err := helpers.FindOldenEraTemplatesDir(false)

	// Assert
	assert.ErrorIs(t, err, common_errors.ErrGameInVDFNotFound)
}

// writeSteamLibraryVDF creates a fake Steam install under homeDir/.local/share/Steam
// whose libraryfolders.vdf lists a single library containing appID; it returns the
// library directory the VDF points at.
func writeSteamLibraryVDF(t *testing.T, homeDir string, appID string) string {
	t.Helper()

	steamDir := filepath.Join(homeDir, ".local", "share", "Steam")
	libraryDir := filepath.Join(steamDir, "library")
	require.NoError(t, os.MkdirAll(filepath.Join(steamDir, "steamapps"), 0o750))

	vdfContent := `"libraryfolders"
{
	"0"
	{
		"path"		"` + libraryDir + `"
		"apps"
		{
			"` + appID + `"		"123456"
		}
	}
}
`
	vdfPath := filepath.Join(steamDir, "steamapps", "libraryfolders.vdf")
	require.NoError(t, os.WriteFile(vdfPath, []byte(vdfContent), 0o600))

	return libraryDir
}
