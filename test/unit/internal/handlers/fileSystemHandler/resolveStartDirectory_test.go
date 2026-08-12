package fileSystemHandler_test

import (
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenStartDirectoryIsRequested_DelegatesToPathResolution(t *testing.T) {
	t.Parallel()
	// Arrange
	handler, mocks := newHandlerWithMocks()
	preferred := gofakeit.LetterN(6)
	mocks.pathResolution.On("ResolveStartDirectory", preferred).Return("resolved")

	// Act
	handler.ResolveStartDirectory(preferred)

	// Assert
	mocks.pathResolution.AssertCalled(t, "ResolveStartDirectory", preferred)
}

func TestWhenStartDirectoryIsResolved_ReturnsItUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	handler, mocks := newHandlerWithMocks()
	mocks.pathResolution.On("ResolveStartDirectory", "wanted").Return("resolved")

	// Act
	resolved := handler.ResolveStartDirectory("wanted")

	// Assert
	assert.Equal(t, "resolved", resolved)
}
