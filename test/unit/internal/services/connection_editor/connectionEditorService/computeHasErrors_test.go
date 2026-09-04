package connectionEditorService_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/connection_editor"
	zone_services "github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
	"github.com/stretchr/testify/assert"
)

func TestWhenEveryConnectionEndpointExists_ReturnsFalse(t *testing.T) {
	t.Parallel()
	// Arrange
	service := connection_editor.NewConnectionEditorService(zone_services.NewZoneTierService())
	zones := []template_model.Zone{{Name: "Spawn-A"}, {Name: "Neutral-1"}}
	connections := []template_model.Connection{{From: "Spawn-A", To: "Neutral-1"}}

	// Act
	hasErrors := service.ComputeHasErrors(zones, connections)

	// Assert
	assert.False(t, hasErrors)
}

func TestWhenFromZoneIsMissing_ReturnsTrue(t *testing.T) {
	t.Parallel()
	// Arrange
	service := connection_editor.NewConnectionEditorService(zone_services.NewZoneTierService())
	zones := []template_model.Zone{{Name: "Neutral-1"}}
	connections := []template_model.Connection{{From: "Spawn-A", To: "Neutral-1"}}

	// Act
	hasErrors := service.ComputeHasErrors(zones, connections)

	// Assert
	assert.True(t, hasErrors)
}

func TestWhenToZoneIsMissing_ReturnsTrue(t *testing.T) {
	t.Parallel()
	// Arrange
	service := connection_editor.NewConnectionEditorService(zone_services.NewZoneTierService())
	zones := []template_model.Zone{{Name: "Spawn-A"}}
	connections := []template_model.Connection{{From: "Spawn-A", To: "Neutral-99"}}

	// Act
	hasErrors := service.ComputeHasErrors(zones, connections)

	// Assert
	assert.True(t, hasErrors)
}
