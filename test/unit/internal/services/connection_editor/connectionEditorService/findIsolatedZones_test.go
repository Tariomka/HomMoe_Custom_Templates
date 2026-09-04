package connectionEditorService_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/connection_editor"
	zone_services "github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
	"github.com/stretchr/testify/assert"
)

func TestWhenZoneHasNoConnections_ReturnsThatZone(t *testing.T) {
	t.Parallel()
	// Arrange
	service := connection_editor.NewConnectionEditorService(zone_services.NewZoneTierService())
	zones := []template_model.Zone{{Name: "Spawn-A"}, {Name: "Neutral-1"}, {Name: "Neutral-2"}}
	connections := []template_model.Connection{{From: "Spawn-A", To: "Neutral-1"}}

	// Act
	isolated := service.FindIsolatedZones(zones, connections)

	// Assert
	assert.Equal(t, []string{"Neutral-2"}, isolated)
}

func TestWhenEveryZoneIsConnected_ReturnsNoZones(t *testing.T) {
	t.Parallel()
	// Arrange
	service := connection_editor.NewConnectionEditorService(zone_services.NewZoneTierService())
	zones := []template_model.Zone{{Name: "Spawn-A"}, {Name: "Neutral-1"}}
	connections := []template_model.Connection{{From: "Spawn-A", To: "Neutral-1"}}

	// Act
	isolated := service.FindIsolatedZones(zones, connections)

	// Assert
	assert.Empty(t, isolated)
}
