package test_helpers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers/handler_interfaces"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/file_system"
)

// NewFileSystemHandler builds the real filesystem handler over the real
// services. Tests that only need the GUI driver to behave as it does in
// production - opening a picker at the working directory, for instance - use
// this instead of FileSystemHandlerMock, which is reserved for tests that
// assert on the interaction itself.
func NewFileSystemHandler() handler_interfaces.IFileSystemHandler {
	return handlers.NewFileSystemHandler(
		file_system.NewDirectoryBrowserService(),
		file_system.NewPathResolutionService())
}
