package interfaces

import "github.com/Tariomka/hommoe_custom_templates/internal/dtos"

type IStatePersistenceHandler interface {
	LoadState(path string, fixIssues bool) (*dtos.EditorStateDto, []string, error)
	SaveState(stateDto dtos.EditorStateSaveDto) (string, error)
}
