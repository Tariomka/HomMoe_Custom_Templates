package fileSystemHandler_test

import (
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenDirectoryCreationIsRequested_DelegatesToTheDirectoryBrowser(t *testing.T) {
	t.Parallel()
	// Arrange
	handler, mocks := newHandlerWithMocks()
	parent := gofakeit.LetterN(6)
	name := gofakeit.LetterN(6)
	mocks.directoryBrowser.On("CreateDirectory", parent, name).Return("created", nil)

	// Act
	_, err := handler.CreateDirectory(parent, name)

	// Assert
	require.NoError(t, err)
	mocks.directoryBrowser.AssertCalled(t, "CreateDirectory", parent, name)
}

func TestWhenDirectoryIsCreated_ReturnsTheCreatedPath(t *testing.T) {
	t.Parallel()
	// Arrange
	handler, mocks := newHandlerWithMocks()
	mocks.directoryBrowser.On("CreateDirectory", "parent", "child").Return("parent/child", nil)

	// Act
	created, err := handler.CreateDirectory("parent", "child")

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "parent/child", created)
}

func TestWhenDirectoryCannotBeCreated_ReturnsTheError(t *testing.T) {
	t.Parallel()
	// Arrange
	handler, mocks := newHandlerWithMocks()
	expectedError := errors.New("denied")
	mocks.directoryBrowser.On("CreateDirectory", "parent", "child").Return("", expectedError)

	// Act
	_, err := handler.CreateDirectory("parent", "child")

	// Assert
	assert.ErrorIs(t, err, expectedError)
}
