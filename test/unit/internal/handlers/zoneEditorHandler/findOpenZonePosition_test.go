package zoneEditorHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenAnOpenPositionIsRequested_ReturnsTheEditorsPosition(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newZoneEditorHandlerFixture()
	occupied := []data.Vec2[float64]{{X: gofakeit.Float64Range(0, 1), Y: gofakeit.Float64Range(0, 1)}}
	expected := data.NewVec2(gofakeit.Float64Range(0, 1), gofakeit.Float64Range(0, 1))
	fixture.zoneEditor.On("FindOpenPosition", occupied).Return(expected)

	// Act
	position := fixture.handler.FindOpenZonePosition(occupied)

	// Assert
	assert.Equal(t, expected, position)
}
