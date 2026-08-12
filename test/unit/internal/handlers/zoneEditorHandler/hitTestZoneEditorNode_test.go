package zoneEditorHandler_test

import (
	"image"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenANodeHitTestIsRequested_ReturnsTheServiceVerdict(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newZoneEditorHandlerFixture()
	expected := gofakeit.Word()
	request := dtos.ZoneEditorHitTestRequestDto{
		Position:   image.Pt(gofakeit.Number(0, 700), gofakeit.Number(0, 700)),
		Positions:  map[string]image.Point{expected: image.Pt(10, 10)},
		ZoneRadius: gofakeit.Number(1, 60),
	}
	fixture.geometry.
		On("HitTestNode", request.Position, request.Positions, request.ZoneRadius).
		Return(expected)

	// Act
	name := fixture.handler.HitTestZoneEditorNode(request)

	// Assert
	assert.Equal(t, expected, name)
}
