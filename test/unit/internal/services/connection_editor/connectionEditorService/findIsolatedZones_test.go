package connectionEditorService_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/connection_editor"
	zone_services "github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
	"github.com/stretchr/testify/assert"
)

func TestWhenZoneHasNoConnections_ReturnsThatZone(t *testing.T) {
	t.Parallel()
	// Arrange
	service := connection_editor.NewConnectionEditorService(zone_services.NewZoneClassifier())
	zones := []entities.Zone{{Name: "Spawn-A"}, {Name: "Neutral-1"}, {Name: "Neutral-2"}}
	connections := []entities.Connection{{From: "Spawn-A", To: "Neutral-1"}}

	// Act
	isolated := service.FindIsolatedZones(zones, connections)

	// Assert
	assert.Equal(t, []string{"Neutral-2"}, isolated)
}

func TestWhenEveryZoneIsConnected_ReturnsNoZones(t *testing.T) {
	t.Parallel()
	// Arrange
	service := connection_editor.NewConnectionEditorService(zone_services.NewZoneClassifier())
	zones := []entities.Zone{{Name: "Spawn-A"}, {Name: "Neutral-1"}}
	connections := []entities.Connection{{From: "Spawn-A", To: "Neutral-1"}}

	// Act
	isolated := service.FindIsolatedZones(zones, connections)

	// Assert
	assert.Empty(t, isolated)
}
