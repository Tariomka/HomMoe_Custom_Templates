package fileSystemHandler_test

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers/handler_interfaces"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
)

type handlerMocks struct {
	directoryBrowser *test_helpers.DirectoryBrowserServiceMock
	pathResolution   *test_helpers.PathResolutionServiceMock
}

// newHandlerWithMocks builds the handler over freshly stubbed services so each
// test owns its own expectations.
func newHandlerWithMocks() (handler_interfaces.IFileSystemHandler, handlerMocks) {
	mocks := handlerMocks{
		directoryBrowser: &test_helpers.DirectoryBrowserServiceMock{},
		pathResolution:   &test_helpers.PathResolutionServiceMock{},
	}

	return handlers.NewFileSystemHandler(mocks.directoryBrowser, mocks.pathResolution), mocks
}
