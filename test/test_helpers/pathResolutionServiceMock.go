package test_helpers

import (
	"github.com/stretchr/testify/mock"
)

// PathResolutionServiceMock is a testify mock of
// file_system.IPathResolutionService, used to unit-test collaborators without
// touching the real filesystem.
type PathResolutionServiceMock struct {
	mock.Mock
}

func (this *PathResolutionServiceMock) ResolveStartDirectory(preferred string) string {
	arguments := this.Called(preferred)
	return arguments.String(0)
}

func (this *PathResolutionServiceMock) ParentDirectory(current string) string {
	arguments := this.Called(current)
	return arguments.String(0)
}

func (this *PathResolutionServiceMock) ResolveSaveTarget(directory, name, requiredSuffix string) (string, bool) {
	arguments := this.Called(directory, name, requiredSuffix)
	return arguments.String(0), arguments.Bool(1)
}

func (this *PathResolutionServiceMock) PathExists(path string) bool {
	arguments := this.Called(path)
	return arguments.Bool(0)
}

func (this *PathResolutionServiceMock) DirectoryExists(path string) bool {
	arguments := this.Called(path)
	return arguments.Bool(0)
}
