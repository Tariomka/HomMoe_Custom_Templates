package zoneEditorGeometryService_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/connection_editor"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
)

func TestWhenTheServiceIsConstructed_ItIsReadyToUse(t *testing.T) {
	t.Parallel()
	// Arrange
	previewLayout := &test_helpers.PreviewLayoutServiceMock{}

	// Act
	service := connection_editor.NewZoneEditorGeometryService(previewLayout)

	// Assert
	assert.NotNil(t, service)
}
