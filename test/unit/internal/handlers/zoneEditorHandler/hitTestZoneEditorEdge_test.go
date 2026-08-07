package zoneEditorHandler_test

import (
	"image"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenAnEdgeHitTestIsRequested_ReturnsTheServiceVerdict(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newZoneEditorHandlerFixture()
	position := image.Pt(gofakeit.Number(0, 700), gofakeit.Number(0, 700))
	edges := []models.ZoneEditorEdge{{ConnectionIndex: gofakeit.Number(0, 9)}}
	expected := gofakeit.Number(0, 9)
	fixture.geometry.On("HitTestEdge", position, edges).Return(expected)

	// Act
	index := fixture.handler.HitTestZoneEditorEdge(position, edges)

	// Assert
	assert.Equal(t, expected, index)
}
