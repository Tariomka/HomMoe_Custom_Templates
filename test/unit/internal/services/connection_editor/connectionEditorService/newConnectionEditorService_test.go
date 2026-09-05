package connectionEditorService_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/connection_editor"
	zone_services "github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
	"github.com/stretchr/testify/assert"
)

func TestWhenServiceIsCreated_ReturnsInstance(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	service := connection_editor.NewConnectionEditorService(zone_services.NewZoneTierService())

	// Assert
	assert.NotNil(t, service)
}
