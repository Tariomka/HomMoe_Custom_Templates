package editor_state_dto

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
)

// EditorStateDto is the serialized .gen.json file produced and consumed by the
// editor. The model is embedded anonymously so the persisted object stays flat
// and every field selector keeps resolving.
type EditorStateDto struct {
	editor_state_model.EditorStateModel
}

func NewDefaultEditorStateDto() EditorStateDto {
	return EditorStateDto{EditorStateModel: editor_state_model.NewDefaultEditorStateModel()}
}

// Clone shadows the promoted EditorStateModel.Clone so callers still holding a
// DTO get a DTO back. Temporary: it goes away once every caller is on the model.
func (this *EditorStateDto) Clone() EditorStateDto {
	return EditorStateDto{EditorStateModel: this.EditorStateModel.Clone()}
}

// LayoutDefiningOptionsChanged shadows the promoted method so it keeps taking a
// DTO. Temporary, as with Clone.
func (this *EditorStateDto) LayoutDefiningOptionsChanged(incoming *EditorStateDto) bool {
	return this.EditorStateModel.LayoutDefiningOptionsChanged(&incoming.EditorStateModel)
}

// DiffCastleSettings shadows the promoted method so it keeps taking a DTO.
// Temporary, as with Clone.
func (this *EditorStateDto) DiffCastleSettings(
	incoming *EditorStateDto,
) editor_state_model.CastleSettingChanges {
	return this.EditorStateModel.DiffCastleSettings(&incoming.EditorStateModel)
}

// EqualsIgnoringManualEdits shadows the promoted method so it keeps taking a
// DTO. Temporary, as with Clone.
func (this *EditorStateDto) EqualsIgnoringManualEdits(other *EditorStateDto) bool {
	return this.EditorStateModel.EqualsIgnoringManualEdits(&other.EditorStateModel)
}
