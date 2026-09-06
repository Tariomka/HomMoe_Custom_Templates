package repositories

import (
	"encoding/json/v2"
	"errors"
	"os"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state/editor_state_v1"
)

// LegacyEditorStateRepository reads .gen.json files written before schema v2,
// whose manualPosition is an array the current entity cannot decode.
type LegacyEditorStateRepository struct{}

func NewLegacyEditorStateRepository() IFileRepository[editor_state_v1.EditorState] {
	return &LegacyEditorStateRepository{}
}

func (this *LegacyEditorStateRepository) Load(filePath string, target *editor_state_v1.EditorState) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, target)
}

// Save always fails: v1 is a frozen read shape, and writing one would produce a
// file this build could no longer round-trip.
func (this *LegacyEditorStateRepository) Save(_, _ string, _ editor_state_v1.EditorState) (string, error) {
	return "", errors.New("the v1 editor state is read-only and exists only so pre-v2 files can be migrated")
}
