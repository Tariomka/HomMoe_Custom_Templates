package connectionEditorService_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/connection_editor"
	zone_services "github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
	"github.com/brianvoe/gofakeit/v7"
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

func TestWhenThereAreNoZones_ReturnsNil(t *testing.T) {
	t.Parallel()
	// Arrange
	service := connection_editor.NewConnectionEditorService(zone_services.NewZoneTierService())
	connections := []template_model.Connection{{From: gofakeit.Word(), To: gofakeit.Word()}}

	// Act
	isolated := service.FindIsolatedZones(nil, connections)

	// Assert
	assert.Nil(t, isolated)
}

func TestWhenThereAreNoConnections_ReturnsEveryZoneInOrder(t *testing.T) {
	t.Parallel()
	// Arrange
	service := connection_editor.NewConnectionEditorService(zone_services.NewZoneTierService())
	zones := []template_model.Zone{{Name: "Neutral-2"}, {Name: "Spawn-A"}, {Name: "Neutral-1"}}

	// Act
	isolated := service.FindIsolatedZones(zones, nil)

	// Assert
	assert.Equal(t, []string{"Neutral-2", "Spawn-A", "Neutral-1"}, isolated)
}

func TestWhenAZoneIsOnlyAConnectionTarget_ItIsNotIsolated(t *testing.T) {
	t.Parallel()
	// Arrange
	service := connection_editor.NewConnectionEditorService(zone_services.NewZoneTierService())
	zones := []template_model.Zone{{Name: "Neutral-1"}}
	connections := []template_model.Connection{{From: "Spawn-A", To: "Neutral-1"}}

	// Act
	isolated := service.FindIsolatedZones(zones, connections)

	// Assert
	assert.Nil(t, isolated)
}

func TestWhenAConnectionJoinsUnknownZones_KnownZonesStayIsolated(t *testing.T) {
	t.Parallel()
	// Arrange
	service := connection_editor.NewConnectionEditorService(zone_services.NewZoneTierService())
	zones := []template_model.Zone{{Name: "Neutral-1"}}
	connections := []template_model.Connection{{From: "Spawn-A", To: "Missing"}}

	// Act
	isolated := service.FindIsolatedZones(zones, connections)

	// Assert
	assert.Equal(t, []string{"Neutral-1"}, isolated)
}

func TestWhenTwoUnreferencedZonesShareAName_ReportsBoth(t *testing.T) {
	t.Parallel()
	// Arrange
	service := connection_editor.NewConnectionEditorService(zone_services.NewZoneTierService())
	zones := []template_model.Zone{{Name: "Neutral-1"}, {Name: "Neutral-1"}}

	// Act
	isolated := service.FindIsolatedZones(zones, nil)

	// Assert
	assert.Equal(t, []string{"Neutral-1", "Neutral-1"}, isolated)
}

func TestWhenAZoneOnlyConnectsToItself_ItIsNotIsolated(t *testing.T) {
	t.Parallel()
	// Arrange
	service := connection_editor.NewConnectionEditorService(zone_services.NewZoneTierService())
	zones := []template_model.Zone{{Name: "Neutral-1"}}
	connections := []template_model.Connection{{From: "Neutral-1", To: "Neutral-1"}}

	// Act
	isolated := service.FindIsolatedZones(zones, connections)

	// Assert
	assert.Nil(t, isolated)
}
