package validators

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
)

// ValidationIssue describes a single invalid EditorStateDto value together
// with the correction that would resolve it. Validation never modifies the
// state; callers decide whether to apply the fix.
type ValidationIssue struct {
	Message string
	fix     func(state *editor_state_dto.EditorStateDto)
}

// Fix applies this issue's correction to the given state.
func (this ValidationIssue) Fix(state *editor_state_dto.EditorStateDto) {
	this.fix(state)
}
