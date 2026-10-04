package validators

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
)

// ValidationIssue describes a single invalid EditorStateModel value together
// with the correction that would resolve it. Validation never modifies the
// state; callers decide whether to apply the fix. A blocking issue has no
// correction: the state must be rejected with its Rejection error.
type ValidationIssue struct {
	Message   string
	fix       func(state *editor_state_model.EditorState)
	rejection error
}

// Fix applies this issue's correction to the given state. It does nothing for
// a blocking issue.
func (this ValidationIssue) Fix(state *editor_state_model.EditorState) {
	if this.IsBlocking() {
		return
	}
	this.fix(state)
}

func (this ValidationIssue) IsBlocking() bool {
	return this.rejection != nil
}

// Rejection returns the error a blocking issue rejects the state with, or nil.
func (this ValidationIssue) Rejection() error {
	return this.rejection
}
