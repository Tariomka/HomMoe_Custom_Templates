package directoryBrowserService_test

import (
	"os"
	"path/filepath"
	"testing"

	internal_constants "github.com/Tariomka/hommoe_custom_templates/internal/common/constants"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/stretchr/testify/require"
)

// makeDirectory creates directory names inside root.
func makeDirectories(t *testing.T, root string, names ...string) {
	t.Helper()
	for _, name := range names {
		require.NoError(t, os.Mkdir(filepath.Join(root, name), internal_constants.FolderPermission))
	}
}

// makeFiles creates empty files names inside root.
func makeFiles(t *testing.T, root string, names ...string) {
	t.Helper()
	for _, name := range names {
		require.NoError(t, os.WriteFile(filepath.Join(root, name), nil, internal_constants.FilePermission))
	}
}

// entryNames projects a listing down to its names so ordering assertions read
// as a single expected slice.
func entryNames(entries []models.DirectoryEntry) []string {
	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		names = append(names, entry.Name)
	}

	return names
}
