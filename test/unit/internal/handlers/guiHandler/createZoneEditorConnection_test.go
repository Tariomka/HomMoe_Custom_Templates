package guiHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers"
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
	expected := entities.Connection{
		From:                 "Spawn-A",
		To:                   "Spawn-B",
		ConnectionType:       "Direct",
		GuardValue:           30000,
		GuardZone:            "Spawn-A",
		GuardMatchGroup:      "rnd_guard_A_B",
		GuardWeeklyIncrement: 0.15,
		IsUserAdded:          true,
	}

	// Act
	result := handler.CreateZoneEditorConnection(request)

	// Assert
	assert.Equal(t, expected, result)
}
