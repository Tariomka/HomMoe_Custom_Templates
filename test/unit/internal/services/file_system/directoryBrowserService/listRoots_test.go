package directoryBrowserService_test

import (
	"runtime"
	"strings"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/file_system"
	"github.com/stretchr/testify/assert"
)

func TestWhenRootsAreListedOnWindows_ReturnsAtLeastOneVolume(t *testing.T) {
	t.Parallel()
	if runtime.GOOS != "windows" {
		t.Skip("volume roots only exist on Windows")
	}
	// Arrange
	service := file_system.NewDirectoryBrowserService()

	// Act
	roots := service.ListRoots()

	// Assert
	assert.NotEmpty(t, roots)
}

func TestWhenRootsAreListedOnWindows_EachRootIsADriveDirectory(t *testing.T) {
	t.Parallel()
	if runtime.GOOS != "windows" {
		t.Skip("volume roots only exist on Windows")
	}
	// Arrange
	service := file_system.NewDirectoryBrowserService()

	// Act
	roots := service.ListRoots()

	// Assert
	for _, root := range roots {
		assert.Equal(t, models.DirectoryEntry{Name: root.Path, Path: root.Path, IsDir: true}, root)
	}
}

func TestWhenRootsAreListedOnWindows_EachNameIsADriveLetterPath(t *testing.T) {
	t.Parallel()
	if runtime.GOOS != "windows" {
		t.Skip("volume roots only exist on Windows")
	}
	// Arrange
	service := file_system.NewDirectoryBrowserService()

	// Act
	roots := service.ListRoots()

	// Assert
	for _, root := range roots {
		assert.True(t, strings.HasSuffix(root.Name, `:\`), "unexpected root name %q", root.Name)
	}
}

func TestWhenRootsAreListedOnNonWindows_ReturnsNothing(t *testing.T) {
	t.Parallel()
	if runtime.GOOS == "windows" {
		t.Skip("drive letters exist on Windows")
	}
	// Arrange
	service := file_system.NewDirectoryBrowserService()

	// Act
	roots := service.ListRoots()

	// Assert
	assert.Empty(t, roots)
}
