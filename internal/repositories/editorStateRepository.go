package repositories

import (
	"encoding/json/v2"
	"os"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
)

const editorStateExtension = ".gen.json"

// EditorStateRepository is the storage boundary for the editor state.
type EditorStateRepository struct {
	writer *atomicFileWriter
}

func NewEditorStateRepository() IFileRepository[editor_state.EditorState] {
	return &EditorStateRepository{writer: newAtomicFileWriter()}
}

func (this *EditorStateRepository) Load(filePath string, target *editor_state.EditorState) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, target)
}

func (this *EditorStateRepository) Save(directory, filename string, entity editor_state.EditorState) (string, error) {
	return this.writer.WriteJSON(directory, filename, editorStateExtension, &entity)
}
