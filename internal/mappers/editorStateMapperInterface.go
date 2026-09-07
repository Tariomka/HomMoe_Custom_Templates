package mappers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
)

type IEditorStateMapper interface {
	NewDefaultEntity() editor_state.EditorState
	ToEntity(state editor_state_model.EditorState) editor_state.EditorState
	ToModel(entity editor_state.EditorState) editor_state_model.EditorState
}
