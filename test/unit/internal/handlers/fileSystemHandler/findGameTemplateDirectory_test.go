package fileSystemHandler_test

import (
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenGameTemplateDirectoryIsRequested_DelegatesToTheDetector(t *testing.T) {
	t.Parallel()
	// Arrange
	handler, mocks := newHandlerWithMocks()
	mocks.pathResolution.On("FindGameTemplateDirectory").Return(gofakeit.LetterN(6), nil)

	// Act
	_, _ = handler.FindGameTemplateDirectory()

	// Assert
	mocks.pathResolution.AssertCalled(t, "FindGameTemplateDirectory")
}

func TestWhenGameTemplateDirectoryIsFound_ReturnsTheDetectedPath(t *testing.T) {
	t.Parallel()
	// Arrange
	handler, mocks := newHandlerWithMocks()
	expectedDirectory := gofakeit.LetterN(6)
	mocks.pathResolution.On("FindGameTemplateDirectory").Return(expectedDirectory, nil)

	// Act
	directory, err := handler.FindGameTemplateDirectory()

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedDirectory, directory)
}

func TestWhenGameTemplateDirectoryLookupFails_PropagatesTheError(t *testing.T) {
	t.Parallel()
	// Arrange
	handler, mocks := newHandlerWithMocks()
	expectedErr := errors.New(gofakeit.Sentence(3))
	mocks.pathResolution.On("FindGameTemplateDirectory").Return("", expectedErr)

	// Act
	_, err := handler.FindGameTemplateDirectory()

	// Assert
	assert.ErrorIs(t, err, expectedErr)
}

// The facade adds no policy of its own, so a detector that reports both a path
// and an error is passed through unchanged; deciding what that combination
// authorizes belongs to the caller.
func TestWhenGameTemplateDirectoryLookupReportsAPathAndAnError_PassesBothThrough(t *testing.T) {
	t.Parallel()
	// Arrange
	handler, mocks := newHandlerWithMocks()
	expectedDirectory := gofakeit.LetterN(6)
	mocks.pathResolution.On("FindGameTemplateDirectory").
		Return(expectedDirectory, errors.New(gofakeit.Sentence(3)))

	// Act
	directory, err := handler.FindGameTemplateDirectory()

	// Assert
	require.Error(t, err)
	assert.Equal(t, expectedDirectory, directory)
}
