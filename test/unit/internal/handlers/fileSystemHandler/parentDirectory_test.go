package fileSystemHandler_test

import (
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenParentDirectoryIsRequested_DelegatesToPathResolution(t *testing.T) {
	t.Parallel()
	// Arrange
	handler, mocks := newHandlerWithMocks()
	current := gofakeit.LetterN(6)
	mocks.pathResolution.On("ParentDirectory", current).Return("parent")

	// Act
	handler.ParentDirectory(current)

	// Assert
	mocks.pathResolution.AssertCalled(t, "ParentDirectory", current)
}

func TestWhenParentDirectoryIsResolved_ReturnsItUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	handler, mocks := newHandlerWithMocks()
	mocks.pathResolution.On("ParentDirectory", "current").Return("parent")

	// Act
	parent := handler.ParentDirectory("current")

	// Assert
	assert.Equal(t, "parent", parent)
}
