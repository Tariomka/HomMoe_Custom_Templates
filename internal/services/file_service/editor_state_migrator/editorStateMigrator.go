package editor_state_migrator

import (
	"encoding/json/v2"
	"os"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state/editor_state_v1"
	"github.com/Tariomka/hommoe_custom_templates/internal/repositories"
)

// firstVec2SchemaVersion is the first version whose manualPosition is an object
// rather than a two-element array. Anything below it needs the legacy
// repository, including v0, which carries no schemaVersion key and probes as 0.
const firstVec2SchemaVersion = 2

// EditorStateMigrator picks the repository a .gen.json can actually be read
// with. The version has to be probed before anything else is decoded:
// encoding/json/v2 hard-fails an array read into a struct, so a pre-v2 file fed
// to the current entity does not degrade, it refuses to load.
type EditorStateMigrator struct {
	editorStateRepository       repositories.IFileRepository[editor_state.EditorState]
	legacyEditorStateRepository repositories.IFileRepository[editor_state_v1.EditorState]
}

func NewEditorStateMigrator(
	editorStateRepository repositories.IFileRepository[editor_state.EditorState],
	legacyEditorStateRepository repositories.IFileRepository[editor_state_v1.EditorState]) IEditorStateMigrator {
	return &EditorStateMigrator{
		editorStateRepository:       editorStateRepository,
		legacyEditorStateRepository: legacyEditorStateRepository,
	}
}

// schemaVersionProbe reads the version key and nothing else. It is declared
// here instead of reusing editor_state.SchemaOptions because the version key is
// the one part of the format that may never move; sharing the current struct
// would let a rename of it break every old file at once.
type schemaVersionProbe struct {
	SchemaVersion int `json:"schemaVersion"`
}

// Load fills target from the file. The caller supplies target already seeded
// with the defaults, so a key the file omits keeps its default instead of
// collapsing to a zero value.
func (this *EditorStateMigrator) Load(filePath string, target *editor_state.EditorState) error {
	version, err := probeSchemaVersion(filePath)
	if err != nil {
		return err
	}

	if version > editor_state.CurrentEditorStateSchemaVersion {
		return newUnsupportedSchemaVersionError(version)
	}

	if version >= firstVec2SchemaVersion {
		return this.editorStateRepository.Load(filePath, target)
	}

	legacy := newV1Seed(*target)
	if err := this.legacyEditorStateRepository.Load(filePath, &legacy); err != nil {
		return err
	}

	*target = MigrateToV2(legacy)
	return nil
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
