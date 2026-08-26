package test_helpers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers/handler_interfaces"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/editor"
)

// NewRegenerationHandler builds the real regeneration handler over the real
// decision service. The service is pure and takes the frame time as an
// argument, so tests that drive the debounce get production behaviour without
// needing a mock.
func NewRegenerationHandler() handler_interfaces.IRegenerationHandler {
	return handlers.NewRegenerationHandler(editor.NewRegenerationDecisionService())
}
