package pathResolutionService_test

import (
	"os"
	"path/filepath"
	"testing"

	internal_constants "github.com/Tariomka/hommoe_custom_templates/internal/common/constants"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/file_system"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenPreferredDirectoryExists_ReturnsIt(t *testing.T) {
	t.Parallel()
	// Arrange
	service := file_system.NewPathResolutionService()
	existing := t.TempDir()

	// Act
	resolved := service.ResolveStartDirectory(existing)

	// Assert
	assert.Equal(t, filepath.Clean(existing), resolved)
}

func TestWhenPreferredDirectoryIsPadded_StillResolvesToIt(t *testing.T) {
	t.Parallel()
	// Arrange
	service := file_system.NewPathResolutionService()
	existing := t.TempDir()

	// Act
	resolved := service.ResolveStartDirectory("  " + existing + "  ")

	// Assert
	assert.Equal(t, filepath.Clean(existing), resolved)
}

func TestWhenPreferredPathIsAFile_ReturnsItsContainingDirectory(t *testing.T) {
	t.Parallel()
	// Arrange
	service := file_system.NewPathResolutionService()
	existing := t.TempDir()
	filePath := filepath.Join(existing, "template.gen.json")
	require.NoError(t, os.WriteFile(filePath, nil, internal_constants.FilePermission))

	// Act
	resolved := service.ResolveStartDirectory(filePath)

	// Assert
	assert.Equal(t, filepath.Clean(existing), resolved)
}

func TestWhenPreferredDirectoryIsMissing_ClimbsToTheNearestExistingAncestor(t *testing.T) {
	t.Parallel()
	// Arrange
	service := file_system.NewPathResolutionService()
	existing := t.TempDir()
	missing := filepath.Join(existing, gofakeit.LetterN(6), gofakeit.LetterN(6))

	// Act
	resolved := service.ResolveStartDirectory(missing)

	// Assert
	assert.Equal(t, filepath.Clean(existing), resolved)
}

func TestWhenPreferredDirectoryIsRelative_ResolvesItAgainstTheWorkingDirectory(t *testing.T) {
	t.Parallel()
	// Arrange
	service := file_system.NewPathResolutionService()
	workingDirectory, err := os.Getwd()
	require.NoError(t, err)

	// Act
	resolved := service.ResolveStartDirectory(".")

	// Assert
	assert.Equal(t, filepath.Clean(workingDirectory), resolved)
}

func TestWhenPreferredDirectoryIsBlank_FallsBackToTheHomeDirectory(t *testing.T) {
	t.Parallel()
	// Arrange
	service := file_system.NewPathResolutionService()
	home, err := os.UserHomeDir()
	require.NoError(t, err)

	// Act
	resolved := service.ResolveStartDirectory("")

	// Assert
	assert.Equal(t, home, resolved)
}
