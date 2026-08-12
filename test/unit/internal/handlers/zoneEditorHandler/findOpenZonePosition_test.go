package zoneEditorHandler_test

import (
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenAnOpenPositionIsRequested_ReturnsTheEditorsPosition(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newZoneEditorHandlerFixture()
	occupied := [][2]float64{{gofakeit.Float64Range(0, 1), gofakeit.Float64Range(0, 1)}}
	expected := [2]float64{gofakeit.Float64Range(0, 1), gofakeit.Float64Range(0, 1)}
	fixture.zoneEditor.On("FindOpenPosition", occupied).Return(expected)

	// Act
	position := fixture.handler.FindOpenZonePosition(occupied)

	// Assert
	assert.Equal(t, expected, position)
}
