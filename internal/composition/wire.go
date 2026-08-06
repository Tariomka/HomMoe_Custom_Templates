//go:build wireinject

// The build tag keeps this stub out of every real build; only the wire CLI
// compiles it, and only to read the graph.

package composition

import (
	"github.com/goforj/wire"

	"github.com/Tariomka/hommoe_custom_templates/internal/handlers/handler_interfaces"
)

func InitializeGuiHandler() handler_interfaces.IGuiHandler {
	wire.Build(GuiHandlerSet)
	return nil
}
