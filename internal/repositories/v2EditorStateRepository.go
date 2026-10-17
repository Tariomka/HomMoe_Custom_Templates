package repositories

import (
	"encoding/json/v2"
	"errors"
	"os"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state/editor_state_v2"
)

type V2EditorStateRepository struct{}

func NewV2EditorStateRepository() IFileRepository[editor_state_v2.EditorState] {
	return &V2EditorStateRepository{}
}

func (this *V2EditorStateRepository) Load(filePath string, target *editor_state_v2.EditorState) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, target)
}

func (this *V2EditorStateRepository) Save(_, _ string, _ editor_state_v2.EditorState) (string, error) {
	return "", errors.New("the v2 editor state is read-only and exists only so v2 files can be migrated")
}
