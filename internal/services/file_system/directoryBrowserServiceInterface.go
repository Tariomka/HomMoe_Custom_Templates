package file_system

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
)

// IDirectoryBrowserService reads directory listings and creates directories. It
// owns every "what does the filesystem contain" decision - hidden-entry
// filtering, suffix filtering and listing order - so callers only render what
// they are handed.
type IDirectoryBrowserService interface {
	ListEntries(directory string, filterSuffixes []string, showHidden bool) ([]models.DirectoryEntry, error)
	ListRoots() []models.DirectoryEntry
	CreateDirectory(parent, name string) (string, error)
}
