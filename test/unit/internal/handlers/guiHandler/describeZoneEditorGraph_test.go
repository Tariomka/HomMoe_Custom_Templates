package guiHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/connection_editor"
	"github.com/stretchr/testify/assert"
)

func TestWhenGraphIsDescribed_ReturnsServiceEquivalentStatus(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := handlers.NewGuiHandler()
	zones := []entities.Zone{{Name: "A"}, {Name: "B"}, {Name: "C"}}
	connections := []entities.Connection{{From: "A", To: "Missing"}}
	expected := dtos.ZoneEditorGraphDto{
		HasErrors:         connection_editor.ComputeHasErrors(zones, connections),
		IsolatedZoneCount: len(connection_editor.FindIsolatedZones(zones, connections)),
	}

	// Act
	result := handler.DescribeZoneEditorGraph(zones, connections)

	// Assert
	assert.Equal(t, expected, result)
}
