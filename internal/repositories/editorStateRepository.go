package repositories

import (
	"encoding/json"
	"os"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
)

const editorStateExtension = ".gen.json"

type EditorStateRepository struct {
	writer *atomicFileWriter
}

func NewEditorStateRepository() IFileRepository[dtos.EditorStateDto] {
	return &EditorStateRepository{writer: newAtomicFileWriter()}
}

func (this *EditorStateRepository) Load(filePath string) (dtos.EditorStateDto, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return dtos.EditorStateDto{}, err
	}

	editorState := dtos.NewDefaultEditorStateDto()
	if err = json.Unmarshal(data, &editorState); err != nil {
		return dtos.EditorStateDto{}, err
	}

	return editorState, nil
}

func (this *EditorStateRepository) Save(
	directory string,
	filename string,
	entity dtos.EditorStateDto) (string, error) {
	return this.writer.WriteJSON(directory, filename, editorStateExtension, &entity)
}
