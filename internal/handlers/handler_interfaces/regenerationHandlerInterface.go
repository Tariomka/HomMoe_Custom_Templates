package handler_interfaces

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
)

// IRegenerationHandler is the facade the GUI uses to decide when the live
// preview regenerates and whether manual zone edits survive that regeneration.
//
// The methods are listed flat rather than embedded from the service interface
// so that app/ never has to import internal/services.
type IRegenerationHandler interface {
	DecideRegeneration(request dtos.RegenerationDecisionRequestDto) dtos.RegenerationDecisionDto
	DecideManualEditReapplication(previous, current *editor_state_model.EditorState) dtos.ManualEditDecisionDto
}
