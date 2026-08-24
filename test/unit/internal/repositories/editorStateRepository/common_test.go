package editorStateRepository_test

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/mappers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/repositories"
)

// newRepository builds the repository with the real entity mapper, because the
// mapping is what decides which fields reach the file.
func newRepository() repositories.IFileRepository[editor_state_model.EditorState] {
	return repositories.NewEditorStateRepository(mappers.NewEditorStateEntityMapper())
}
