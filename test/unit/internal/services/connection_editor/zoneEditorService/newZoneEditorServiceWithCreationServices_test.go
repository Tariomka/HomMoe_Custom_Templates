package zoneEditorService_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/connection_editor"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
	"github.com/stretchr/testify/assert"
)

func TestWhenCreationServicesAreProvided_ReturnsUsableService(t *testing.T) {
	t.Parallel()
	// Arrange
	creationServices := zones.NewCreationServices(nil, nil)

	// Act
	service := connection_editor.NewZoneEditorServiceWithCreationServices(creationServices)

	// Assert
	assert.NotNil(t, service)
}
