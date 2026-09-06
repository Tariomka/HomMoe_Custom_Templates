package editor_state_migrator

import (
	"fmt"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
)

// UnsupportedSchemaVersionError reports a .gen.json written by a newer build
// than this one. Loading it best-effort would drop the keys this build does not
// know about, and the next save would write that loss back to disk, so the load
// is refused instead.
type UnsupportedSchemaVersionError struct {
	FileVersion      int
	SupportedVersion int
}

func (this *UnsupportedSchemaVersionError) Error() string {
	return fmt.Sprintf(
		"this file was saved by a newer version of the editor "+
			"(schema version %d; this build understands up to %d) - update the app to open it",
		this.FileVersion,
		this.SupportedVersion)
}

func newUnsupportedSchemaVersionError(fileVersion int) error {
	return &UnsupportedSchemaVersionError{
		FileVersion:      fileVersion,
		SupportedVersion: editor_state.CurrentEditorStateSchemaVersion,
	}
}
