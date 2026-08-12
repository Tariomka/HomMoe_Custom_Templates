package pathResolutionService_test

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/file_system"
	"github.com/stretchr/testify/assert"
)

func TestWhenCurrentDirectoryIsTheSyntheticRootListing_StaysThere(t *testing.T) {
	t.Parallel()
	// Arrange
	service := file_system.NewPathResolutionService()

	// Act
	parent := service.ParentDirectory("")

	// Assert
	assert.Empty(t, parent)
}

func TestWhenCurrentDirectoryIsNested_ReturnsItsParent(t *testing.T) {
	t.Parallel()
	// Arrange
	service := file_system.NewPathResolutionService()
	nested := filepath.Join(t.TempDir(), "child")

	// Act
	parent := service.ParentDirectory(nested)

	// Assert
	assert.Equal(t, filepath.Dir(nested), parent)
}

func TestWhenCurrentDirectoryIsAVolumeRootOnWindows_AscendsToTheRootListing(t *testing.T) {
	t.Parallel()
	if runtime.GOOS != "windows" {
		t.Skip("volume roots only exist on Windows")
	}
	// Arrange
	service := file_system.NewPathResolutionService()

	// Act
	parent := service.ParentDirectory(`C:\`)

	// Assert
	assert.Empty(t, parent)
}

func TestWhenCurrentDirectoryIsTheFilesystemRootOnNonWindows_StaysThere(t *testing.T) {
	t.Parallel()
	if runtime.GOOS == "windows" {
		t.Skip("the single filesystem root only exists off Windows")
	}
	// Arrange
	service := file_system.NewPathResolutionService()

	// Act
	parent := service.ParentDirectory("/")

	// Assert
	assert.Equal(t, "/", parent)
}
