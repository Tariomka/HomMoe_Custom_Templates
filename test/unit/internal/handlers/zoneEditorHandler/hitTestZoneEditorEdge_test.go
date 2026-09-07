package zoneEditorHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenAnEdgeHitTestIsRequested_ReturnsTheServiceVerdict(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newZoneEditorHandlerFixture()
	position := data.NewVec2(gofakeit.Float64Range(0, 700), gofakeit.Float64Range(0, 700))
	edges := []models.ZoneEditorEdge{{ConnectionIndex: gofakeit.Number(0, 9)}}
	expected := gofakeit.Number(0, 9)
	fixture.geometry.On("HitTestEdge", position, edges).Return(expected)

	// Act
	index := fixture.handler.HitTestZoneEditorEdge(position, edges)

	// Assert
	assert.Equal(t, expected, index)
}
