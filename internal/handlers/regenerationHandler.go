package handlers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers/handler_interfaces"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/editor"
)

// regenerationHandler exposes the regeneration decision service as the facade
// the GUI is allowed to depend on. It adds no logic of its own; every method
// delegates, so the policy stays testable inside internal/services.
type regenerationHandler struct {
	regenerationDecision editor.IRegenerationDecisionService
}

func NewRegenerationHandler(
	regenerationDecision editor.IRegenerationDecisionService,
) handler_interfaces.IRegenerationHandler {
	return &regenerationHandler{regenerationDecision: regenerationDecision}
}

func (this *regenerationHandler) DecideRegeneration(
	request dtos.RegenerationDecisionRequestDto,
) dtos.RegenerationDecisionDto {
	return this.regenerationDecision.DecideRegeneration(request)
}

func (this *regenerationHandler) DecideManualEditReapplication(
	previous, current *dtos.EditorStateDto,
) dtos.ManualEditDecisionDto {
	return this.regenerationDecision.DecideManualEditReapplication(previous, current)
}
