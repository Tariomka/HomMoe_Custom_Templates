package test_helpers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/stretchr/testify/mock"
)

// FileSystemHandlerMock is a testify mock of
// handler_interfaces.IFileSystemHandler, used to unit-test collaborators
// without touching the real filesystem.
type FileSystemHandlerMock struct {
	mock.Mock
}

func (this *FileSystemHandlerMock) ListEntries(
	directory string,
	filterSuffixes []string,
	showHidden bool) ([]models.DirectoryEntry, error) {
	arguments := this.Called(directory, filterSuffixes, showHidden)
	entries, _ := arguments.Get(0).([]models.DirectoryEntry)
	return entries, arguments.Error(1)
}

func (this *FileSystemHandlerMock) ListRoots() []models.DirectoryEntry {
	arguments := this.Called()
	roots, _ := arguments.Get(0).([]models.DirectoryEntry)
	return roots
}

func (this *FileSystemHandlerMock) CreateDirectory(parent, name string) (string, error) {
	arguments := this.Called(parent, name)
	return arguments.String(0), arguments.Error(1)
}

func (this *FileSystemHandlerMock) ResolveStartDirectory(preferred string) string {
	arguments := this.Called(preferred)
	return arguments.String(0)
}

func (this *FileSystemHandlerMock) ParentDirectory(current string) string {
	arguments := this.Called(current)
	return arguments.String(0)
}

func (this *FileSystemHandlerMock) ResolveSaveTarget(directory, name, requiredSuffix string) (string, bool) {
	arguments := this.Called(directory, name, requiredSuffix)
	return arguments.String(0), arguments.Bool(1)
}

func (this *FileSystemHandlerMock) PathExists(path string) bool {
	arguments := this.Called(path)
	return arguments.Bool(0)
}
