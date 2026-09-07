package validators

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
)

// ValidationIssue describes a single invalid EditorStateModel value together
// with the correction that would resolve it. Validation never modifies the
// state; callers decide whether to apply the fix.
type ValidationIssue struct {
	Message string
	fix     func(state *editor_state_model.EditorState)
}

// Fix applies this issue's correction to the given state.
func (this ValidationIssue) Fix(state *editor_state_model.EditorState) {
	this.fix(state)
}
