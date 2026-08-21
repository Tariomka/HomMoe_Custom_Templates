package dtos

import (
	"time"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
)

// RegenerationDecisionRequestDto is everything the regeneration decision
// depends on. Previous and Next are nil when absent, which is what
// distinguishes a first generation from a steady state.
//
// Now is the current frame time and DebounceDueAt is when the armed debounce
// window elapses; both are supplied by the caller so the decision stays
// deterministic and testable.
type RegenerationDecisionRequestDto struct {
	Previous      *editor_state_dto.EditorStateDto
	Current       *editor_state_dto.EditorStateDto
	Next          *editor_state_dto.EditorStateDto
	Now           time.Time
	DebounceDueAt time.Time
}
