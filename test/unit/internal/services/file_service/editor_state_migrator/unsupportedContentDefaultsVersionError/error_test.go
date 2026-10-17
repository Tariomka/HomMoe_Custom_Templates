package unsupportedContentDefaultsVersionError_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/file_service/editor_state_migrator"
	"github.com/stretchr/testify/assert"
)

func TestWhenTheDefaultsAreTooNew_TheMessageNamesBothVersions(t *testing.T) {
	t.Parallel()
	// Arrange
	err := &editor_state_migrator.UnsupportedContentDefaultsVersionError{FileVersion: 3, SupportedVersion: 1}

	// Act
	message := err.Error()

	// Assert
	assert.Equal(
		t,
		"this file was saved by a newer version of the editor "+
			"(content defaults version 3; this build understands up to 1) - update the app to open it",
		message)
}
