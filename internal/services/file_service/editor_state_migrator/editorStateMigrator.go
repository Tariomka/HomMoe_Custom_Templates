package editor_state_migrator

import (
	"encoding/json/v2"
	"os"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state/editor_state_v1"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state/editor_state_v2"
	"github.com/Tariomka/hommoe_custom_templates/internal/repositories"
)

const firstVec2SchemaVersion = 2

type EditorStateMigrator struct {
	editorStateRepository   repositories.IFileRepository[editor_state.EditorState]
	v2EditorStateRepository repositories.IFileRepository[editor_state_v2.EditorState]
	v1EditorStateRepository repositories.IFileRepository[editor_state_v1.EditorState]
}

func NewEditorStateMigrator(
	editorStateRepository repositories.IFileRepository[editor_state.EditorState],
	v2EditorStateRepository repositories.IFileRepository[editor_state_v2.EditorState],
	v1EditorStateRepository repositories.IFileRepository[editor_state_v1.EditorState]) IEditorStateMigrator {
	return &EditorStateMigrator{
		editorStateRepository:   editorStateRepository,
		v2EditorStateRepository: v2EditorStateRepository,
		v1EditorStateRepository: v1EditorStateRepository,
	}
}

type schemaVersionProbe struct {
	SchemaVersion int `json:"schemaVersion"`
}

func (this *EditorStateMigrator) Load(filePath string, target *editor_state.EditorState) ([]string, error) {
	version, err := probeSchemaVersion(filePath)
	if err != nil {
		return nil, err
	}

	if version > editor_state.CurrentEditorStateSchemaVersion {
		return nil, newUnsupportedSchemaVersionError(version)
	}

	if version == editor_state.CurrentEditorStateSchemaVersion {
		return nil, this.loadCurrent(filePath, target)
	}

	legacy, err := this.loadLegacy(filePath, version, newV2Seed(*target))
	if err != nil {
		return nil, err
	}

	migrated, warnings := MigrateToV3(legacy)
	*target = migrated
	return warnings, nil
}

func (this *EditorStateMigrator) loadCurrent(filePath string, target *editor_state.EditorState) error {
	loaded := *target
	if err := this.editorStateRepository.Load(filePath, &loaded); err != nil {
		return err
	}

	if loaded.Content.ContentDefaultsVersion > editor_state.CurrentContentDefaultsVersion {
		return newUnsupportedContentDefaultsVersionError(loaded.Content.ContentDefaultsVersion)
	}

	*target = loaded
	return nil
}

// loadLegacy decodes a pre-v3 file into the frozen v2 shape, going through the
// frozen v1 shape first for files whose manualPosition is still an array.
func (this *EditorStateMigrator) loadLegacy(
	filePath string,
	version int,
	seed editor_state_v2.EditorState) (editor_state_v2.EditorState, error) {
	if version >= firstVec2SchemaVersion {
		err := this.v2EditorStateRepository.Load(filePath, &seed)
		return seed, err
	}

	legacy := newV1Seed(seed)
	if err := this.v1EditorStateRepository.Load(filePath, &legacy); err != nil {
		return editor_state_v2.EditorState{}, err
	}

	return MigrateToV2(legacy), nil
}

func probeSchemaVersion(filePath string) (int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	var probe schemaVersionProbe
	if err := json.UnmarshalRead(file, &probe); err != nil {
		return 0, err
	}

	return probe.SchemaVersion, nil
}
