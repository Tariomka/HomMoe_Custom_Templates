package test_helpers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/stretchr/testify/mock"
)

// DirectoryBrowserServiceMock is a testify mock of
// file_system.IDirectoryBrowserService, used to unit-test collaborators without
// touching the real filesystem.
type DirectoryBrowserServiceMock struct {
	mock.Mock
}

func (this *DirectoryBrowserServiceMock) ListEntries(
	directory string,
	filterSuffixes []string,
	showHidden bool) ([]models.DirectoryEntry, error) {
	arguments := this.Called(directory, filterSuffixes, showHidden)
	entries, _ := arguments.Get(0).([]models.DirectoryEntry)
	return entries, arguments.Error(1)
}

func (this *DirectoryBrowserServiceMock) ListRoots() []models.DirectoryEntry {
	arguments := this.Called()
	roots, _ := arguments.Get(0).([]models.DirectoryEntry)
	return roots
}

func (this *DirectoryBrowserServiceMock) CreateDirectory(parent, name string) (string, error) {
	arguments := this.Called(parent, name)
	return arguments.String(0), arguments.Error(1)
}
