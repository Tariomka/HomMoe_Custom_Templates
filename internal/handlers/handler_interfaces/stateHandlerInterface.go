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
	LoadState(path string, fixIssues bool) (*editor_state_dto.EditorStateDto, []string, error)
	SaveState(stateDto editor_state_dto.EditorStateSaveDto) (string, error)
}

// IStateValidationHandler deliberately keeps speaking the Model. It sits on the
// per-frame path - app/gui/models.EditorState.UpdateCurrentState calls it on
// every panel write - so it is the one crossing exempted from the DTO rule.
type IStateValidationHandler interface {
	ValidateEditorState(
		state editor_state_model.EditorState,
		fixIssues bool) editor_state_dto.EditorStateValidationDto
}
