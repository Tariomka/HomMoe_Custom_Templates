package editor

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
)

// IRegenerationDecisionService owns when the live preview regenerates and
// whether hand-made zone edits survive that regeneration.
//
// Both methods are pure: every input arrives as an argument (including the
// current time) and neither mutates the snapshots it is handed. Callers apply
// the returned NextStateAction themselves.
type IRegenerationDecisionService interface {
	DecideRegeneration(request dtos.RegenerationDecisionRequestDto) dtos.RegenerationDecisionDto
	DecideManualEditReapplication(previous, current *editor_state_model.EditorState) dtos.ManualEditDecisionDto
}
