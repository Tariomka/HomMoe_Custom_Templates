package mappers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
)

// EditorStateEntityMapper converts between the persisted editor state and the
// runtime one. The model embeds the entity, so the conversion carries the whole
// group set across in one assignment and only the schema version needs deciding.
type EditorStateEntityMapper struct{}

func NewEditorStateEntityMapper() IEditorStateEntityMapper {
	return &EditorStateEntityMapper{}
}

// NewDefaultEntity is the seed a load decodes over: a key the file omits keeps
// the default instead of collapsing to a zero value. A zero-seeded decode
// cannot tell an absent key from an explicit false or 0, so the seed has to be
// in place before the file is read.
func (this *EditorStateEntityMapper) NewDefaultEntity() editor_state.EditorState {
	return this.ToEntity(editor_state_model.NewDefaultEditorStateModel())
}

func (this *EditorStateEntityMapper) ToEntity(state editor_state_model.EditorState) editor_state.EditorState {
	entity := state.EditorState
	entity.SchemaVersion = editor_state.CurrentEditorStateSchemaVersion
	return entity
}

func (this *EditorStateEntityMapper) ToModel(entity editor_state.EditorState) editor_state_model.EditorState {
	entity.SchemaVersion = editor_state.CurrentEditorStateSchemaVersion
	return editor_state_model.EditorState{EditorState: entity}
}
