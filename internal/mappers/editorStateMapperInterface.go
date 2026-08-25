package mappers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
)

// IEditorStateMapper converts the editor state between the crossing contract
// (DTO) and the runtime type (Model). It is reachable from app/ because the
// GUI maps its stored Model into a request DTO at every handler call site.
type IEditorStateMapper interface {
	ToDto(state editor_state_model.EditorState) editor_state_dto.EditorStateDto
	ToModel(dto editor_state_dto.EditorStateDto) editor_state_model.EditorState

	ToDtoPointer(state *editor_state_model.EditorState) *editor_state_dto.EditorStateDto
	ToModelPointer(dto *editor_state_dto.EditorStateDto) *editor_state_model.EditorState

	ToCastleSettingChangesDto(
		changes editor_state_model.CastleSettingChanges) editor_state_dto.CastleSettingChangesDto
	ToCastleSettingChangesModel(
		dto editor_state_dto.CastleSettingChangesDto) editor_state_model.CastleSettingChanges
}
