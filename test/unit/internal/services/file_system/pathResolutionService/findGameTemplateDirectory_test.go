package pathResolutionService_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/file_system"
	"github.com/stretchr/testify/assert"
)

// The method is a one-line adapter over the platform detector, so the only
// thing worth asserting is that it reports exactly what the detector does on
// this host - installed or not. The detector's own discovery rules are covered
// by its platform fixture tests.
func TestWhenGameTemplateDirectoryIsRequested_ReportsTheDetectedPath(t *testing.T) {
	t.Parallel()
	// Arrange
	service := file_system.NewPathResolutionService()
	expectedDirectory, _ := helpers.FindOldenEraTemplatesDir(false)

	// Act
	directory, _ := service.FindGameTemplateDirectory()

	// Assert
	assert.Equal(t, expectedDirectory, directory)
}

func TestWhenGameTemplateDirectoryIsRequested_ReportsTheDetectorError(t *testing.T) {
	t.Parallel()
	// Arrange
	service := file_system.NewPathResolutionService()
	_, expectedErr := helpers.FindOldenEraTemplatesDir(false)

	// Act
	_, err := service.FindGameTemplateDirectory()

	// Assert
	assert.Equal(t, expectedErr, err)
}
