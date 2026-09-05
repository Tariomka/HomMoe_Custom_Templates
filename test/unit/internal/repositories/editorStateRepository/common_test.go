package editorStateRepository_test

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/Tariomka/hommoe_custom_templates/internal/repositories"
)

func newRepository() repositories.IFileRepository[editor_state.EditorState] {
	return repositories.NewEditorStateRepository()
}
