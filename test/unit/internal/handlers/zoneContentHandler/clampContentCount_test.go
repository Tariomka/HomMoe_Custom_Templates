package zoneContentHandler_test

import (
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenTheContentCountIsClamped_ReturnsTheServiceValue(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newZoneContentHandlerFixture()
	count := gofakeit.Number(1, 100)
	maxCount := gofakeit.Number(1, 100)
	expected := gofakeit.Number(1, 100)
	fixture.contentEditor.On("ClampContentCount", count, maxCount).Return(expected)

	// Act
	clamped := fixture.handler.ClampContentCount(count, maxCount)

	// Assert
	assert.Equal(t, expected, clamped)
}
