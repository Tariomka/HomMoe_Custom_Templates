package pathResolutionService_test

import (
	"os"
	"path/filepath"
	"testing"

	internal_constants "github.com/Tariomka/hommoe_custom_templates/internal/common/constants"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/file_system"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenDirectoryPathIsAnExistingDirectory_ReportsItIsADirectory(t *testing.T) {
	t.Parallel()
	// Arrange
	service := file_system.NewPathResolutionService()

	// Act
	isDirectory := service.DirectoryExists(t.TempDir())

	// Assert
	assert.True(t, isDirectory)
}

func TestWhenDirectoryPathIsAnExistingFile_ReportsItIsNotADirectory(t *testing.T) {
	t.Parallel()
	// Arrange
	service := file_system.NewPathResolutionService()
	filePath := filepath.Join(t.TempDir(), "template.gen.json")
	require.NoError(t, os.WriteFile(filePath, nil, internal_constants.FilePermission))

	// Act
	isDirectory := service.DirectoryExists(filePath)

	// Assert
	assert.False(t, isDirectory)
}

func TestWhenDirectoryPathPointsAtNothing_ReportsItIsNotADirectory(t *testing.T) {
	t.Parallel()
	// Arrange
	service := file_system.NewPathResolutionService()

	// Act
	isDirectory := service.DirectoryExists(filepath.Join(t.TempDir(), "absent"))

	// Assert
	assert.False(t, isDirectory)
}
