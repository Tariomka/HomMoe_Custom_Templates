package guiHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/connection_editor"
	"github.com/stretchr/testify/assert"
)

func TestWhenZoneIsRemoved_ReturnsServiceEquivalentMutation(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := handlers.NewGuiHandler()
	request := dtos.ZoneEditorRemoveRequestDto{
		Zones: []entities.Zone{{Name: "Spawn-A"}, {Name: "Neutral-C"}, {Name: "Neutral-D"}},
		Connections: []entities.Connection{
			{From: "Spawn-A", To: "Neutral-C"},
			{From: "Spawn-A", To: "Neutral-D"},
		},
		ZoneName: "Neutral-C",
	}
	zones, connections := connection_editor.NewZoneEditorService().
		RemoveZone(request.Zones, request.Connections, request.ZoneName)
	expected := dtos.ZoneEditorMutationDto{Zones: zones, Connections: connections}

	// Act
	result := handler.RemoveZoneEditorZone(request)

	// Assert
	assert.Equal(t, expected, result)
}
