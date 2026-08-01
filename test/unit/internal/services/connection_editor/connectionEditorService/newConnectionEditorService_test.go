package connectionEditorService_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/connection_editor"
	zone_services "github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
	"github.com/stretchr/testify/assert"
)

func TestWhenServiceIsCreated_ReturnsInstance(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	service := connection_editor.NewConnectionEditorService(zone_services.NewZoneClassifier())

	// Assert
	assert.NotNil(t, service)
}

func TestWhenClassifierIsNil_ReturnsUsableService(t *testing.T) {
	t.Parallel()
	// Arrange
	service := connection_editor.NewConnectionEditorService(nil)
	zones := []entities.Zone{{Name: "Spawn-A"}, {Name: "Spawn-B"}}
	playerZoneNames := map[string]bool{"Spawn-A": true, "Spawn-B": true}

	// Act
	connection := service.NewDefaultConnection("Spawn-A", "Spawn-B", zones, playerZoneNames)

	// Assert
	assert.Equal(t, 30000, connection.GuardValue)
}
