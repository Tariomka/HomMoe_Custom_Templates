package guiHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/stretchr/testify/assert"
)

func TestWhenPortalIsPicked_ReturnsTheTypedConnectionWithItsFlagKept(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	request := dtos.ZoneEditorConnectionTypeRequestDto{
		Connection:     template_model.Connection{From: "Spawn-A", To: "Spawn-B", ConnectionType: "Direct"},
		ConnectionType: "Portal",
		GenerateRoads:  true,
	}
	expected := template_model.Connection{From: "Spawn-A", To: "Spawn-B", ConnectionType: "Portal"}

	// Act
	result := handler.ChangeZoneEditorConnectionType(request)

	// Assert
	assert.Equal(t, expected, result)
}

func TestWhenDirectIsPickedWithRoadsOff_ReturnsTheRoadlessConnection(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	request := dtos.ZoneEditorConnectionTypeRequestDto{
		Connection:     template_model.Connection{From: "Spawn-A", To: "Spawn-B", ConnectionType: "Portal"},
		ConnectionType: "Direct",
	}
	expected := template_model.Connection{
		From:           "Spawn-A",
		To:             "Spawn-B",
		ConnectionType: "Direct",
		Road:           new(false),
	}

	// Act
	result := handler.ChangeZoneEditorConnectionType(request)

	// Assert
	assert.Equal(t, expected, result)
}
