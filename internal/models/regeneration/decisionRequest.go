package regeneration

import (
	"time"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
)

// DecisionRequest is everything the regeneration decision depends on. Previous
// and Next are nil when absent, which is what distinguishes a first generation
// from a steady state.
//
// Now is the current frame time and DebounceDueAt is when the armed debounce
// window elapses; both are supplied by the caller so the decision stays
// deterministic and testable.
type DecisionRequest struct {
	Previous      *editor_state_model.EditorState
	Current       *editor_state_model.EditorState
	Next          *editor_state_model.EditorState
	Now           time.Time
	DebounceDueAt time.Time
}
