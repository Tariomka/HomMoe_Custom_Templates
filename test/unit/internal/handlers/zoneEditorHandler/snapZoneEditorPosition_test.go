package zoneEditorHandler_test

import (
	"image"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenASnapIsRequested_ReturnsTheServiceResult(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newZoneEditorHandlerFixture()
	request := dtos.ZoneEditorSnapRequestDto{
		Position:    image.Pt(gofakeit.Number(0, 700), gofakeit.Number(0, 700)),
		Positions:   map[string]image.Point{gofakeit.Word(): image.Pt(350, 350)},
		ZoneRadius:  gofakeit.Number(1, 60),
		DraggedZone: gofakeit.Word(),
	}
	expected := models.ZoneEditorSnapResult{
		Position:  image.Pt(gofakeit.Number(0, 700), gofakeit.Number(0, 700)),
		HasGuideY: true,
	}
	fixture.geometry.
		On("SnapPosition", request.Position, request.Positions, request.ZoneRadius, request.DraggedZone).
		Return(expected)

	// Act
	result := fixture.handler.SnapZoneEditorPosition(request)

	// Assert
	assert.Equal(t, expected, result)
}
