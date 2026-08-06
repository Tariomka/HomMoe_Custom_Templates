package validators

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
)

type IEditorStateValidator interface {
	Validate(state *dtos.EditorStateDto) []ValidationIssue
}
