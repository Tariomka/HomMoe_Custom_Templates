package unsupportedSchemaVersionError_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/file_service/editor_state_migrator"
	"github.com/stretchr/testify/assert"
)

// The message is what the user sees behind a failed load, so it has to name the
// file's version and say what to do about it rather than read as a parse error.
func TestWhenTheFileIsTooNew_TheMessageNamesBothVersions(t *testing.T) {
	t.Parallel()
	// Arrange
	err := &editor_state_migrator.UnsupportedSchemaVersionError{FileVersion: 7, SupportedVersion: 2}

	// Act
	message := err.Error()

	// Assert
	assert.Equal(
		t,
		"this file was saved by a newer version of the editor "+
			"(schema version 7; this build understands up to 2) - update the app to open it",
		message)
}
