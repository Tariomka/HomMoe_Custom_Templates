package legacyEditorStateRepository_test

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state/editor_state_v1"
	"github.com/Tariomka/hommoe_custom_templates/internal/repositories"
)

func newRepository() repositories.IFileRepository[editor_state_v1.EditorState] {
	return repositories.NewLegacyEditorStateRepository()
}
