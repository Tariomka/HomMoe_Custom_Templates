package editor_state_dto

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
)

type EditorStateSaveDto struct {
	State      *editor_state_model.EditorState
	OutputPath string
}
