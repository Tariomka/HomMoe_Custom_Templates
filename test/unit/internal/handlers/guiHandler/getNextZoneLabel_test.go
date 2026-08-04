package guiHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/connection_editor"
	"github.com/stretchr/testify/assert"
)

func TestWhenLabelsAreOccupied_ReturnsServiceEquivalentLabel(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	zones := []entities.Zone{{Name: "Spawn-A"}, {Name: "Neutral-B"}}
	expected := connection_editor.NewDefaultZoneEditorService().NextFreeZoneLabel(zones)

	// Act
	result := handler.GetNextZoneLabel(zones)

	// Assert
	assert.Equal(t, expected, result)
}
