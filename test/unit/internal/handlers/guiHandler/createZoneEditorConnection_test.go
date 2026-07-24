package guiHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/connection_editor"
	"github.com/stretchr/testify/assert"
)

func TestWhenConnectionIsCreated_ReturnsServiceEquivalentConnection(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := handlers.NewGuiHandler()
	request := dtos.ZoneEditorConnectionRequestDto{
		From:            "Spawn-A",
		To:              "Spawn-B",
		Zones:           []entities.Zone{{Name: "Spawn-A"}, {Name: "Spawn-B"}},
		PlayerZoneNames: map[string]bool{"Spawn-A": true, "Spawn-B": true},
	}
	expected := connection_editor.NewDefaultConnection(
		request.From, request.To, request.Zones, request.PlayerZoneNames)

	// Act
	result := handler.CreateZoneEditorConnection(request)

	// Assert
	assert.Equal(t, expected, result)
}
