package fileSystemHandler_test

import (
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenExistenceIsChecked_DelegatesToPathResolution(t *testing.T) {
	t.Parallel()
	// Arrange
	handler, mocks := newHandlerWithMocks()
	path := gofakeit.LetterN(6)
	mocks.pathResolution.On("PathExists", path).Return(true)

	// Act
	handler.PathExists(path)

	// Assert
	mocks.pathResolution.AssertCalled(t, "PathExists", path)
}

func TestWhenPathIsReportedPresent_ReturnsTrue(t *testing.T) {
	t.Parallel()
	// Arrange
	handler, mocks := newHandlerWithMocks()
	mocks.pathResolution.On("PathExists", "there").Return(true)

	// Act
	exists := handler.PathExists("there")

	// Assert
	assert.True(t, exists)
}

func TestWhenPathIsReportedAbsent_ReturnsFalse(t *testing.T) {
	t.Parallel()
	// Arrange
	handler, mocks := newHandlerWithMocks()
	mocks.pathResolution.On("PathExists", "gone").Return(false)

	// Act
	exists := handler.PathExists("gone")

	// Assert
	assert.False(t, exists)
}
