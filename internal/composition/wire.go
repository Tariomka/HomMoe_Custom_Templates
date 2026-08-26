//go:build wireinject

// The build tag keeps this stub out of every real build; only the wire CLI
// compiles it, and only to read the graph.

package composition

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers/handler_interfaces"
	"github.com/goforj/wire"
)

func InitializeGuiHandler() handler_interfaces.IGuiHandler {
	wire.Build(GuiHandlerSet)
	return nil
}

func InitializeFileSystemHandler() handler_interfaces.IFileSystemHandler {
	wire.Build(FileSystemSet)
	return nil
}

func InitializeRegenerationHandler() handler_interfaces.IRegenerationHandler {
	wire.Build(RegenerationSet)
	return nil
}
