package handlers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers/handler_interfaces"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/regeneration"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/editor"
)

// regenerationHandler exposes the regeneration decision service as the facade
// the GUI is allowed to depend on. It adds no policy of its own; it only maps
// the crossing DTOs onto the model types the service works in, so the policy
// stays testable inside internal/services.
type regenerationHandler struct {
	regenerationDecision editor.IRegenerationDecisionService
}

func NewRegenerationHandler(
	regenerationDecision editor.IRegenerationDecisionService) handler_interfaces.IRegenerationHandler {
	return &regenerationHandler{
		regenerationDecision: regenerationDecision,
	}
}

func (this *regenerationHandler) DecideRegeneration(
	request dtos.RegenerationDecisionRequestDto) dtos.RegenerationDecisionDto {
	decision := this.regenerationDecision.DecideRegeneration(regeneration.DecisionRequest{
		Previous:      toEditorStateModel(request.Previous),
		Current:       toEditorStateModel(request.Current),
		Next:          toEditorStateModel(request.Next),
		Now:           request.Now,
		DebounceDueAt: request.DebounceDueAt,
	})

	return dtos.RegenerationDecisionDto(decision)
}

func (this *regenerationHandler) DecideManualEditReapplication(
	previous, current *editor_state_dto.EditorStateDto) dtos.ManualEditDecisionDto {
	decision := this.regenerationDecision.DecideManualEditReapplication(
		toEditorStateModel(previous),
		toEditorStateModel(current))

	if decision.ReapplyWithCastleChanges == nil {
		return dtos.ManualEditDecisionDto{}
	}

	return dtos.ManualEditDecisionDto{ReapplyWithCastleChanges: decision.ReapplyWithCastleChanges}
}

func toEditorStateModel(dto *editor_state_dto.EditorStateDto) *editor_state_model.EditorState {
	if dto == nil {
		return nil
	}

	return &dto.EditorState
}
