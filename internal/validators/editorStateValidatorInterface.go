package validators

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
)

type IEditorStateValidator interface {
	Validate(state *editor_state_model.EditorState) []ValidationIssue
}
