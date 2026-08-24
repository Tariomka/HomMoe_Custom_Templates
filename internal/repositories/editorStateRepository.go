package repositories

import (
	"encoding/json"
	"os"

	"github.com/Tariomka/hommoe_custom_templates/internal/mappers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
)

const editorStateExtension = ".gen.json"

// EditorStateRepository is the storage boundary for the editor state: it speaks
// the entity to the file and the model to everything above it.
type EditorStateRepository struct {
	mapper mappers.IEditorStateEntityMapper
	writer *atomicFileWriter
}

func NewEditorStateRepository(
	mapper mappers.IEditorStateEntityMapper) IFileRepository[editor_state_model.EditorState] {
	return &EditorStateRepository{mapper: mapper, writer: newAtomicFileWriter()}
}

func (this *EditorStateRepository) Load(filePath string) (editor_state_model.EditorState, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return editor_state_model.EditorState{}, err
	}

	// Decoding over the defaults, rather than over a zero value, is what lets a
	// setting the file omits keep its default.
	entity := this.mapper.ToEntity(editor_state_model.NewDefaultEditorStateModel())
	if err = json.Unmarshal(data, &entity); err != nil {
		return editor_state_model.EditorState{}, err
	}

	return this.mapper.ToModel(entity), nil
}

func (this *EditorStateRepository) Save(
	directory string,
	filename string,
	state editor_state_model.EditorState) (string, error) {
	entity := this.mapper.ToEntity(state)

	return this.writer.WriteJSON(directory, filename, editorStateExtension, &entity)
}
