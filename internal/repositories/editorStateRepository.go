package repositories

import "github.com/Tariomka/hommoe_custom_templates/internal/dtos"

type EditorStateRepository struct{}

func NewEditorStateRepository() IFileRepository[dtos.EditorStateDto] {
	return &EditorStateRepository{}
}

func (this *EditorStateRepository) Load(filePath string) (dtos.EditorStateDto, error) {
	return dtos.EditorStateDto{}, nil
}

func (this *EditorStateRepository) Save(directory string, filename string, entity dtos.EditorStateDto) (string, error) {
	return "", nil
}
