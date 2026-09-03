package zoneEditorHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenNextZoneLabelIsRequested_ReturnsTheEditorsLabel(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newZoneEditorHandlerFixture()
	zones := []template_model.Zone{{Name: gofakeit.Word()}}
	expected := gofakeit.Word()
	fixture.zoneEditor.On("NextFreeZoneLabel", zones).Return(expected)

	// Act
	label := fixture.handler.GetNextZoneLabel(zones)

	// Assert
	assert.Equal(t, expected, label)
}
