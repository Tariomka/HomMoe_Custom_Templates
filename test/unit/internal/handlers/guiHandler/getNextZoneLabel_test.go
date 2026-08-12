package guiHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
)

func TestWhenLabelsAreOccupied_ReturnsServiceEquivalentLabel(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	zones := []entities.Zone{{Name: "Spawn-A"}, {Name: "Neutral-B"}}
	expected := test_helpers.NewZoneEditorService().NextFreeZoneLabel(zones)

	// Act
	result := handler.GetNextZoneLabel(zones)

	// Assert
	assert.Equal(t, expected, result)
}
