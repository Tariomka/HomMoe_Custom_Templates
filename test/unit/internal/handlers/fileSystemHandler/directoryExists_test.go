package fileSystemHandler_test

import (
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenDirectoryExistenceIsChecked_DelegatesToPathResolution(t *testing.T) {
	t.Parallel()
	// Arrange
	handler, mocks := newHandlerWithMocks()
	path := gofakeit.LetterN(6)
	mocks.pathResolution.On("DirectoryExists", path).Return(true)

	// Act
	handler.DirectoryExists(path)

	// Assert
	mocks.pathResolution.AssertCalled(t, "DirectoryExists", path)
}

func TestWhenPathIsReportedToBeADirectory_ReturnsTrue(t *testing.T) {
	t.Parallel()
	// Arrange
	handler, mocks := newHandlerWithMocks()
	mocks.pathResolution.On("DirectoryExists", "there").Return(true)

	// Act
	isDirectory := handler.DirectoryExists("there")

	// Assert
	assert.True(t, isDirectory)
}

func TestWhenPathIsReportedNotToBeADirectory_ReturnsFalse(t *testing.T) {
	t.Parallel()
	// Arrange
	handler, mocks := newHandlerWithMocks()
	mocks.pathResolution.On("DirectoryExists", "there").Return(false)

	// Act
	isDirectory := handler.DirectoryExists("there")

	// Assert
	assert.False(t, isDirectory)
}
