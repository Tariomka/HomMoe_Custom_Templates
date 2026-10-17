package editor_state_migrator

import (
	"fmt"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
)

type UnsupportedContentDefaultsVersionError struct {
	FileVersion      int
	SupportedVersion int
}

func (this *UnsupportedContentDefaultsVersionError) Error() string {
	return fmt.Sprintf(
		"this file was saved by a newer version of the editor "+
			"(content defaults version %d; this build understands up to %d) - update the app to open it",
		this.FileVersion,
		this.SupportedVersion)
}

func newUnsupportedContentDefaultsVersionError(fileVersion int) error {
	return &UnsupportedContentDefaultsVersionError{
		FileVersion:      fileVersion,
		SupportedVersion: editor_state.CurrentContentDefaultsVersion,
	}
}
