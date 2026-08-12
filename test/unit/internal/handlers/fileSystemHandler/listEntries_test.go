package fileSystemHandler_test

import (
	"errors"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenEntriesAreRequested_DelegatesToTheDirectoryBrowser(t *testing.T) {
	t.Parallel()
	// Arrange
	handler, mocks := newHandlerWithMocks()
	directory := gofakeit.LetterN(6)
	suffixes := []string{".gen.json"}
	mocks.directoryBrowser.On("ListEntries", directory, suffixes, true).Return([]models.DirectoryEntry{}, nil)

	// Act
	_, err := handler.ListEntries(directory, suffixes, true)

	// Assert
	require.NoError(t, err)
	mocks.directoryBrowser.AssertCalled(t, "ListEntries", directory, suffixes, true)
}

func TestWhenEntriesAreListed_ReturnsThemUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	handler, mocks := newHandlerWithMocks()
	expected := []models.DirectoryEntry{{Name: "nested", Path: "nested", IsDir: true}}
	mocks.directoryBrowser.On("ListEntries", "any", []string(nil), false).Return(expected, nil)

	// Act
	actual, err := handler.ListEntries("any", nil, false)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestWhenEntriesCannotBeListed_ReturnsTheError(t *testing.T) {
	t.Parallel()
	// Arrange
	handler, mocks := newHandlerWithMocks()
	expectedError := errors.New("unreadable")
	mocks.directoryBrowser.On("ListEntries", "any", []string(nil), false).
		Return([]models.DirectoryEntry(nil), expectedError)

	// Act
	_, err := handler.ListEntries("any", nil, false)

	// Assert
	assert.ErrorIs(t, err, expectedError)
}
