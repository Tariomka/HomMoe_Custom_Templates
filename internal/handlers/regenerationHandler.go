package handlers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers/handler_interfaces"
	"github.com/Tariomka/hommoe_custom_templates/internal/mappers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/regeneration"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/editor"
)

// regenerationHandler exposes the regeneration decision service as the facade
// the GUI is allowed to depend on. It adds no policy of its own; it only maps
// the crossing DTOs onto the model types the service works in, so the policy
// stays testable inside internal/services.
type regenerationHandler struct {
	regenerationDecision editor.IRegenerationDecisionService
	editorStateMapper    mappers.IEditorStateMapper
}

func NewRegenerationHandler(
	regenerationDecision editor.IRegenerationDecisionService,
	editorStateMapper mappers.IEditorStateMapper) handler_interfaces.IRegenerationHandler {
	return &regenerationHandler{
		regenerationDecision: regenerationDecision,
		editorStateMapper:    editorStateMapper,
	}
}

func (this *regenerationHandler) DecideRegeneration(
	request dtos.RegenerationDecisionRequestDto) dtos.RegenerationDecisionDto {
	decision := this.regenerationDecision.DecideRegeneration(regeneration.DecisionRequest{
		Previous:      this.editorStateMapper.ToModelPointer(request.Previous),
		Current:       this.editorStateMapper.ToModelPointer(request.Current),
		Next:          this.editorStateMapper.ToModelPointer(request.Next),
		Now:           request.Now,
		DebounceDueAt: request.DebounceDueAt,
	})

	return dtos.RegenerationDecisionDto(decision)
}

func (this *regenerationHandler) DecideManualEditReapplication(
	previous, current *editor_state_dto.EditorStateDto) dtos.ManualEditDecisionDto {
	decision := this.regenerationDecision.DecideManualEditReapplication(
		this.editorStateMapper.ToModelPointer(previous),
		this.editorStateMapper.ToModelPointer(current))

	if decision.ReapplyWithCastleChanges == nil {
		return dtos.ManualEditDecisionDto{}
	}

	changes := this.editorStateMapper.ToCastleSettingChangesDto(*decision.ReapplyWithCastleChanges)
	return dtos.ManualEditDecisionDto{ReapplyWithCastleChanges: &changes}
}
