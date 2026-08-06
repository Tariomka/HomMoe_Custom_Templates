package directoryBrowserService_test

import (
	"path/filepath"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/file_system"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenDirectoryCannotBeRead_ReturnsError(t *testing.T) {
	t.Parallel()
	// Arrange
	service := file_system.NewDirectoryBrowserService()
	missing := filepath.Join(t.TempDir(), "does-not-exist")

	// Act
	_, err := service.ListEntries(missing, nil, true)

	// Assert
	assert.Error(t, err)
}

func TestWhenDirectoryCannotBeRead_ReturnsNoEntries(t *testing.T) {
	t.Parallel()
	// Arrange
	service := file_system.NewDirectoryBrowserService()
	missing := filepath.Join(t.TempDir(), "does-not-exist")

	// Act
	entries, _ := service.ListEntries(missing, nil, true)

	// Assert
	assert.Empty(t, entries)
}

func TestWhenDirectoryHoldsFilesAndFolders_ListsFoldersFirst(t *testing.T) {
	t.Parallel()
	// Arrange
	service := file_system.NewDirectoryBrowserService()
	root := t.TempDir()
	makeFiles(t, root, "alpha.txt", "zulu.txt")
	makeDirectories(t, root, "beta", "yankee")

	// Act
	entries, err := service.ListEntries(root, nil, true)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, []string{"beta", "yankee", "alpha.txt", "zulu.txt"}, entryNames(entries))
}

func TestWhenNamesDifferOnlyInCase_SortsCaseInsensitively(t *testing.T) {
	t.Parallel()
	// Arrange
	service := file_system.NewDirectoryBrowserService()
	root := t.TempDir()
	makeFiles(t, root, "Banana.txt", "apple.txt", "Cherry.txt")

	// Act
	entries, err := service.ListEntries(root, nil, true)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, []string{"apple.txt", "Banana.txt", "Cherry.txt"}, entryNames(entries))
}

func TestWhenHiddenEntriesAreSuppressed_OmitsDotPrefixedEntries(t *testing.T) {
	t.Parallel()
	// Arrange
	service := file_system.NewDirectoryBrowserService()
	root := t.TempDir()
	makeFiles(t, root, ".secret", "visible.txt")

	// Act
	entries, err := service.ListEntries(root, nil, false)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, []string{"visible.txt"}, entryNames(entries))
}

func TestWhenHiddenEntriesAreShown_IncludesDotPrefixedEntries(t *testing.T) {
	t.Parallel()
	// Arrange
	service := file_system.NewDirectoryBrowserService()
	root := t.TempDir()
	makeFiles(t, root, ".secret", "visible.txt")

	// Act
	entries, err := service.ListEntries(root, nil, true)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, []string{".secret", "visible.txt"}, entryNames(entries))
}

func TestWhenFilterSuffixesAreGiven_OmitsNonMatchingFiles(t *testing.T) {
	t.Parallel()
	// Arrange
	service := file_system.NewDirectoryBrowserService()
	root := t.TempDir()
	makeFiles(t, root, "keep.gen.json", "drop.txt")

	// Act
	entries, err := service.ListEntries(root, []string{".gen.json"}, true)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, []string{"keep.gen.json"}, entryNames(entries))
}

func TestWhenFilterSuffixCaseDiffers_StillMatchesTheFile(t *testing.T) {
	t.Parallel()
	// Arrange
	service := file_system.NewDirectoryBrowserService()
	root := t.TempDir()
	makeFiles(t, root, "Keep.GEN.Json")

	// Act
	entries, err := service.ListEntries(root, []string{".gen.json"}, true)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, []string{"Keep.GEN.Json"}, entryNames(entries))
}

func TestWhenFilterSuffixesAreGiven_StillListsSubdirectories(t *testing.T) {
	t.Parallel()
	// Arrange
	service := file_system.NewDirectoryBrowserService()
	root := t.TempDir()
	makeDirectories(t, root, "nested")

	// Act
	entries, err := service.ListEntries(root, []string{".gen.json"}, true)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, []string{"nested"}, entryNames(entries))
}

func TestWhenFilterSuffixesAreEmpty_ListsEveryFile(t *testing.T) {
	t.Parallel()
	// Arrange
	service := file_system.NewDirectoryBrowserService()
	root := t.TempDir()
	makeFiles(t, root, "one.txt", "two.json")

	// Act
	entries, err := service.ListEntries(root, nil, true)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, []string{"one.txt", "two.json"}, entryNames(entries))
}

func TestWhenEntryIsListed_CarriesItsFullPathAndKind(t *testing.T) {
	t.Parallel()
	// Arrange
	service := file_system.NewDirectoryBrowserService()
	root := t.TempDir()
	makeDirectories(t, root, "nested")

	// Act
	entries, err := service.ListEntries(root, nil, true)

	// Assert
	require.NoError(t, err)
	assert.Equal(
		t,
		[]models.DirectoryEntry{{Name: "nested", Path: filepath.Join(root, "nested"), IsDir: true}},
		entries)
}

func TestWhenDirectoryIsEmpty_ReturnsEmptyListing(t *testing.T) {
	t.Parallel()
	// Arrange
	service := file_system.NewDirectoryBrowserService()

	// Act
	entries, err := service.ListEntries(t.TempDir(), nil, true)

	// Assert
	require.NoError(t, err)
	assert.Empty(t, entries)
}
