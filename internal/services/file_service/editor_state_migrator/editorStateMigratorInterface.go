package editor_state_migrator

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
)

type IEditorStateMigrator interface {
	Load(filePath string, target *editor_state.EditorState) ([]string, error)
}
