package guiHandler_test

import (
	"image"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/stretchr/testify/assert"
)

func TestWhenASnapRequested_HoldsTheDraggedZoneOnItsNeighboursGuide(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	request := dtos.ZoneEditorSnapRequestDto{
		Position:    image.Pt(200, 355),
		Positions:   map[string]image.Point{"Spawn-A": image.Pt(350, 350)},
		ZoneRadius:  38,
		DraggedZone: "Spawn-B",
	}

	// Act
	result := handler.SnapZoneEditorPosition(request)

	// Assert
	assert.Equal(t, image.Pt(201, 350), result.Position)
}
