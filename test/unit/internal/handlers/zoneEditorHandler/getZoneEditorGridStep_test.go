package zoneEditorHandler_test

import (
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenTheGridStepIsRequested_ReturnsTheServiceValue(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newZoneEditorHandlerFixture()
	zoneRadius := gofakeit.Float64Range(1, 60)
	expected := gofakeit.Float64Range(1, 40)
	fixture.geometry.On("GridStep", zoneRadius).Return(expected)

	// Act
	step := fixture.handler.GetZoneEditorGridStep(zoneRadius)

	// Assert
	assert.InDelta(t, expected, step, 1e-9)
}
