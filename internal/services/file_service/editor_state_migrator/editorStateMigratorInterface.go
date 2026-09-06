package editor_state_migrator

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
)

// IEditorStateMigrator loads a .gen.json into the current entity shape,
// whichever schema version wrote it.
type IEditorStateMigrator interface {
	Load(filePath string, target *editor_state.EditorState) error
}
