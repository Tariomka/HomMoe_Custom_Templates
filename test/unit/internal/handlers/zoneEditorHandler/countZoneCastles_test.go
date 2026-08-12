package zoneEditorHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenZoneCastlesAreCounted_ReturnsTheEditorsCount(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newZoneEditorHandlerFixture()
	zone := entities.Zone{Name: gofakeit.Word()}
	expected := gofakeit.IntRange(0, 4)
	fixture.zoneEditor.On("CountZoneCastles", zone).Return(expected)

	// Act
	count := fixture.handler.CountZoneCastles(zone)

	// Assert
	assert.Equal(t, expected, count)
}
