package zoneContentHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenContentItemsAreSorted_ReturnsTheServiceOrdering(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newZoneContentHandlerFixture()
	items := []models.SidMapping{{Name: gofakeit.Word()}}
	expected := []models.SidMapping{{Name: gofakeit.Sentence(2)}}
	fixture.contentEditor.On("SortContentItemsByName", items).Return(expected)

	// Act
	sorted := fixture.handler.SortContentItemsByName(items)

	// Assert
	assert.Equal(t, expected, sorted)
}
