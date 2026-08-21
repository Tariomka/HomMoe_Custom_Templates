package validators

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
)

type IEditorStateValidator interface {
	Validate(state *editor_state_dto.EditorStateDto) []ValidationIssue
}
