package handler_interfaces

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
)

type IStateHandler interface {
	IStateValidationHandler
	IStatePersistenceHandler
}

type IStatePersistenceHandler interface {
	LoadState(path string, fixIssues bool) (*editor_state_model.EditorState, []string, error)
	SaveState(stateDto editor_state_dto.EditorStateSaveDto) (string, error)
}

type IStateValidationHandler interface {
	ValidateEditorState(
		state editor_state_model.EditorState,
		fixIssues bool) editor_state_dto.EditorStateValidationDto
}
