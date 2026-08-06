package handlers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers/handler_interfaces"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/file_system"
)

// fileSystemHandler exposes the two filesystem services as the single facade
// the GUI is allowed to depend on. It adds no logic of its own; every method
// delegates, so the policy stays testable inside internal/services.
type fileSystemHandler struct {
	directoryBrowser file_system.IDirectoryBrowserService
	pathResolution   file_system.IPathResolutionService
}

func NewFileSystemHandler(
	directoryBrowser file_system.IDirectoryBrowserService,
	pathResolution file_system.IPathResolutionService) handler_interfaces.IFileSystemHandler {
	return &fileSystemHandler{
		directoryBrowser: directoryBrowser,
		pathResolution:   pathResolution,
	}
}

func (this *fileSystemHandler) ListEntries(
	directory string,
	filterSuffixes []string,
	showHidden bool) ([]models.DirectoryEntry, error) {
	return this.directoryBrowser.ListEntries(directory, filterSuffixes, showHidden)
}

func (this *fileSystemHandler) ListRoots() []models.DirectoryEntry {
	return this.directoryBrowser.ListRoots()
}

func (this *fileSystemHandler) CreateDirectory(parent, name string) (string, error) {
	return this.directoryBrowser.CreateDirectory(parent, name)
}

func (this *fileSystemHandler) ResolveStartDirectory(preferred string) string {
	return this.pathResolution.ResolveStartDirectory(preferred)
}

func (this *fileSystemHandler) ParentDirectory(current string) string {
	return this.pathResolution.ParentDirectory(current)
}

func (this *fileSystemHandler) ResolveSaveTarget(directory, name, requiredSuffix string) (string, bool) {
	return this.pathResolution.ResolveSaveTarget(directory, name, requiredSuffix)
}

func (this *fileSystemHandler) PathExists(path string) bool {
	return this.pathResolution.PathExists(path)
}
