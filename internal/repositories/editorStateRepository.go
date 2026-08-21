package repositories

import (
	"encoding/json"
	"os"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
)

const editorStateExtension = ".gen.json"

type EditorStateRepository struct {
	writer *atomicFileWriter
}

func NewEditorStateRepository() IFileRepository[editor_state_dto.EditorStateDto] {
	return &EditorStateRepository{writer: newAtomicFileWriter()}
}

func (this *EditorStateRepository) Load(filePath string) (editor_state_dto.EditorStateDto, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return editor_state_dto.EditorStateDto{}, err
	}

	editorState := editor_state_dto.NewDefaultEditorStateDto()
	if err = json.Unmarshal(data, &editorState); err != nil {
		return editor_state_dto.EditorStateDto{}, err
	}

	return editorState, nil
}

func (this *EditorStateRepository) Save(
	directory string,
	filename string,
	entity editor_state_dto.EditorStateDto) (string, error) {
	return this.writer.WriteJSON(directory, filename, editorStateExtension, &entity)
}
