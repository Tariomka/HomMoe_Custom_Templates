package editor_state_dto

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
)

type EditorStateDto struct {
	editor_state_model.EditorState
}

func NewDefaultEditorStateDto() EditorStateDto {
	return EditorStateDto{EditorState: editor_state_model.NewDefaultEditorStateModel()}
}
