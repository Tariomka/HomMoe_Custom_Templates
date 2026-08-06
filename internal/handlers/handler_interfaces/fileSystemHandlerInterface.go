package handler_interfaces

import "github.com/Tariomka/hommoe_custom_templates/internal/models"

// IFileSystemHandler is the GUI's only door to the local filesystem. It is a
// standalone seam rather than part of IGuiHandler because browsing the disk is
// unrelated to editing a template: the file explorer needs it, and nothing that
// needs IGuiHandler needs this.
type IFileSystemHandler interface {
	ListEntries(directory string, filterSuffixes []string, showHidden bool) ([]models.DirectoryEntry, error)
	ListRoots() []models.DirectoryEntry
	CreateDirectory(parent, name string) (string, error)
	ResolveStartDirectory(preferred string) string
	ParentDirectory(current string) string
	ResolveSaveTarget(directory, name, requiredSuffix string) (string, bool)
	PathExists(path string) bool
}
