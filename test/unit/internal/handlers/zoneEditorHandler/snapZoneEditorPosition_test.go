package zoneEditorHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenASnapIsRequested_ReturnsTheServiceResult(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newZoneEditorHandlerFixture()
	request := dtos.ZoneEditorSnapRequestDto{
		Position:    data.NewVec2(gofakeit.Float64Range(0, 700), gofakeit.Float64Range(0, 700)),
		Positions:   map[string]models.Position{gofakeit.Word(): data.NewVec2(350.0, 350.0)},
		ZoneRadius:  gofakeit.Float64Range(1, 60),
		DraggedZone: gofakeit.Word(),
	}
	expected := models.ZoneEditorSnapResult{
		Position:  data.NewVec2(gofakeit.Float64Range(0, 700), gofakeit.Float64Range(0, 700)),
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
