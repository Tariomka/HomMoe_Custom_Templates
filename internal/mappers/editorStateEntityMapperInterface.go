package mappers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
)

type IEditorStateEntityMapper interface {
	ToEntity(state editor_state_model.EditorState) editor_state.EditorStateEntity
	ToModel(entity editor_state.EditorStateEntity) editor_state_model.EditorState
}
