package interfaces

import "github.com/Tariomka/hommoe_custom_templates/internal/dtos"

type IStateValidationHandler interface {
	ValidateEditorState(state dtos.EditorStateDto, fixIssues bool) dtos.EditorStateValidationDto
}
