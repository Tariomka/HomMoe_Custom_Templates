package handler_interfaces

import "github.com/Tariomka/hommoe_custom_templates/internal/dtos"

type IStateHandler interface {
	IStateValidationHandler
	IStatePersistenceHandler
}

type IStatePersistenceHandler interface {
	LoadState(path string, fixIssues bool) (*dtos.EditorStateDto, []string, error)
	SaveState(stateDto dtos.EditorStateSaveDto) (string, error)
}

type IStateValidationHandler interface {
	ValidateEditorState(state dtos.EditorStateDto, fixIssues bool) dtos.EditorStateValidationDto
}
