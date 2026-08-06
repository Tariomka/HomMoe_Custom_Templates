package directoryBrowserService_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_errors"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/file_system"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenFolderNameIsBlank_ReturnsEmptyNameError(t *testing.T) {
	t.Parallel()
	// Arrange
	service := file_system.NewDirectoryBrowserService()

	// Act
	_, err := service.CreateDirectory(t.TempDir(), "   ")

	// Assert
	assert.ErrorIs(t, err, common_errors.ErrDirectoryNameEmpty)
}

func TestWhenFolderNameIsARelativeOrSeparatedPath_ReturnsInvalidNameError(t *testing.T) {
	t.Parallel()
	invalidNames := map[string]string{
		"CurrentDirectoryToken":  ".",
		"ParentDirectoryToken":   "..",
		"ContainsForwardSlash":   "nested/child",
		"ContainsPathSeparator":  "nested" + string(os.PathSeparator) + "child",
		"ContainsVolumeSeparato": "C:evil",
	}
	for scenario, invalidName := range invalidNames {
		t.Run(scenario+"_ReturnsInvalidNameError", func(t *testing.T) {
			t.Parallel()
			// Arrange
			service := file_system.NewDirectoryBrowserService()

			// Act
			_, err := service.CreateDirectory(t.TempDir(), invalidName)

			// Assert
			assert.ErrorIs(t, err, common_errors.ErrDirectoryNameInvalid)
		})
	}
}

func TestWhenFolderNameIsValid_CreatesTheDirectory(t *testing.T) {
	t.Parallel()
	// Arrange
	service := file_system.NewDirectoryBrowserService()
	parent := t.TempDir()
	name := gofakeit.LetterN(8)

	// Act
	_, err := service.CreateDirectory(parent, name)

	// Assert
	require.NoError(t, err)
	assert.DirExists(t, filepath.Join(parent, name))
}

func TestWhenFolderIsCreated_ReturnsTheCreatedPath(t *testing.T) {
	t.Parallel()
	// Arrange
	service := file_system.NewDirectoryBrowserService()
	parent := t.TempDir()
	name := gofakeit.LetterN(8)

	// Act
	created, err := service.CreateDirectory(parent, name)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(parent, name), created)
}

func TestWhenFolderNameIsPadded_CreatesTheTrimmedName(t *testing.T) {
	t.Parallel()
	// Arrange
	service := file_system.NewDirectoryBrowserService()
	parent := t.TempDir()
	name := gofakeit.LetterN(8)

	// Act
	created, err := service.CreateDirectory(parent, "  "+name+"  ")

	// Assert
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(parent, name), created)
}

func TestWhenFolderAlreadyExists_ReturnsTheUnderlyingError(t *testing.T) {
	t.Parallel()
	// Arrange
	service := file_system.NewDirectoryBrowserService()
	parent := t.TempDir()
	name := gofakeit.LetterN(8)
	_, err := service.CreateDirectory(parent, name)
	require.NoError(t, err)

	// Act
	_, err = service.CreateDirectory(parent, name)

	// Assert
	assert.ErrorIs(t, err, os.ErrExist)
}

func TestWhenParentDoesNotExist_ReturnsTheUnderlyingError(t *testing.T) {
	t.Parallel()
	// Arrange
	service := file_system.NewDirectoryBrowserService()
	missingParent := filepath.Join(t.TempDir(), "absent")

	// Act
	_, err := service.CreateDirectory(missingParent, gofakeit.LetterN(8))

	// Assert
	assert.ErrorIs(t, err, os.ErrNotExist)
}
