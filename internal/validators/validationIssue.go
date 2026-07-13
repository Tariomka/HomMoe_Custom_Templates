package validators

import "github.com/Tariomka/hommoe_custom_templates/internal/dtos"

// ValidationIssue describes a single invalid EditorStateDto value together
// with the correction that would resolve it. Validation never modifies the
// state; callers decide whether to apply the fix.
type ValidationIssue struct {
	Message string
	fix     func(state *dtos.EditorStateDto)
}

// Fix applies this issue's correction to the given state.
func (this ValidationIssue) Fix(state *dtos.EditorStateDto) {
	this.fix(state)
}
