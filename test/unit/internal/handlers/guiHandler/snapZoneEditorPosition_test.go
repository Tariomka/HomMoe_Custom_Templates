package guiHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestWhenASnapRequested_HoldsTheDraggedZoneOnItsNeighboursGuide(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	request := dtos.ZoneEditorSnapRequestDto{
		Position:    data.NewVec2(200.0, 355.0),
		Positions:   map[string]models.Position{"Spawn-A": data.NewVec2(350.0, 350.0)},
		ZoneRadius:  38,
		DraggedZone: "Spawn-B",
	}

	// Act
	result := handler.SnapZoneEditorPosition(request)

	// Assert
	assert.Equal(t, data.NewVec2(200.85714285714286, 350.0), result.Position)
}
