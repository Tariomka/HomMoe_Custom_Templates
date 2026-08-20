package zoneEditorHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenANodeHitTestIsRequested_ReturnsTheServiceVerdict(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newZoneEditorHandlerFixture()
	expected := gofakeit.Word()
	request := dtos.ZoneEditorHitTestRequestDto{
		Position:   data.NewVec2(gofakeit.Float64Range(0, 700), gofakeit.Float64Range(0, 700)),
		Positions:  map[string]models.Position{expected: data.NewVec2(10.0, 10.0)},
		ZoneRadius: gofakeit.Float64Range(1, 60),
	}
	fixture.geometry.
		On("HitTestNode", request.Position, request.Positions, request.ZoneRadius).
		Return(expected)

	// Act
	name := fixture.handler.HitTestZoneEditorNode(request)

	// Assert
	assert.Equal(t, expected, name)
}
