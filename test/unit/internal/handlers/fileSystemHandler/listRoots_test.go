package fileSystemHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestWhenRootsAreRequested_DelegatesToTheDirectoryBrowser(t *testing.T) {
	t.Parallel()
	// Arrange
	handler, mocks := newHandlerWithMocks()
	mocks.directoryBrowser.On("ListRoots").Return([]models.DirectoryEntry{})

	// Act
	handler.ListRoots()

	// Assert
	mocks.directoryBrowser.AssertCalled(t, "ListRoots")
}

func TestWhenRootsAreListed_ReturnsThemUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	handler, mocks := newHandlerWithMocks()
	expected := []models.DirectoryEntry{{Name: `C:\`, Path: `C:\`, IsDir: true}}
	mocks.directoryBrowser.On("ListRoots").Return(expected)

	// Act
	actual := handler.ListRoots()

	// Assert
	assert.Equal(t, expected, actual)
}
