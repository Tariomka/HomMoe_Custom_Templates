package guiHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
)

func TestWhenZoneIsRemoved_ReturnsServiceEquivalentMutation(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	request := dtos.ZoneEditorRemoveRequestDto{
		Zones: []template_model.Zone{{Name: "Spawn-A"}, {Name: "Neutral-C"}, {Name: "Neutral-D"}},
		Connections: []template_model.Connection{
			{From: "Spawn-A", To: "Neutral-C"},
			{From: "Spawn-A", To: "Neutral-D"},
		},
		ZoneName: "Neutral-C",
	}
	zones, connections := test_helpers.NewZoneEditorService().
		RemoveZone(request.Zones, request.Connections, request.ZoneName)
	expected := dtos.ZoneEditorMutationDto{Zones: zones, Connections: connections}

	// Act
	result := handler.RemoveZoneEditorZone(request)

	// Assert
	assert.Equal(t, expected, result)
}
