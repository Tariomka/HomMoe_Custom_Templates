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

func TestWhenPathPointsAtAnExistingFile_ReportsItExists(t *testing.T) {
	t.Parallel()
	// Arrange
	service := file_system.NewPathResolutionService()
	filePath := filepath.Join(t.TempDir(), "template.gen.json")
	require.NoError(t, os.WriteFile(filePath, nil, internal_constants.FilePermission))

	// Act
	exists := service.PathExists(filePath)

	// Assert
	assert.True(t, exists)
}

func TestWhenPathPointsAtAnExistingDirectory_ReportsItExists(t *testing.T) {
	t.Parallel()
	// Arrange
	service := file_system.NewPathResolutionService()

	// Act
	exists := service.PathExists(t.TempDir())

	// Assert
	assert.True(t, exists)
}

func TestWhenPathPointsAtNothing_ReportsItIsAbsent(t *testing.T) {
	t.Parallel()
	// Arrange
	service := file_system.NewPathResolutionService()

	// Act
	exists := service.PathExists(filepath.Join(t.TempDir(), "absent.gen.json"))

	// Assert
	assert.False(t, exists)
}
