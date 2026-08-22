package editor_state_dto

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
)

// EditorStateDto is the serialized .gen.json file produced and consumed by the
// editor. The model is embedded anonymously so the persisted object stays flat
// and every field selector keeps resolving.
type EditorStateDto struct {
	editor_state_model.EditorState
}

func NewDefaultEditorStateDto() EditorStateDto {
	return EditorStateDto{EditorState: editor_state_model.NewDefaultEditorStateModel()}
}

// NewEditorStateDto wraps a runtime model for persistence. It and Model are the
// only two places that know how the DTO stores the model, so reshaping the DTO
// stays confined to this file.
func NewEditorStateDto(state editor_state_model.EditorState) EditorStateDto {
	return EditorStateDto{EditorState: state}
}

// Model returns the runtime model carried by the DTO. The result aliases the
// receiver's slices; callers that keep it past the receiver's lifetime clone.
func (this *EditorStateDto) Model() *editor_state_model.EditorState {
	return &this.EditorState
}
