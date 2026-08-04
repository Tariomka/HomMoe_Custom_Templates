package guiHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/connection_editor"
	"github.com/stretchr/testify/assert"
)

func TestWhenZoneContainsMixedMainObjects_ReturnsServiceEquivalentCount(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	zone := entities.Zone{MainObjects: []entities.MainObject{
		{Type: "Spawn"},
		{Type: "City"},
		{Type: "AbandonedOutpost"},
		{Type: "City"},
	}}
	expected := connection_editor.NewDefaultZoneEditorService().CountZoneCastles(zone)

	// Act
	result := handler.CountZoneCastles(zone)

	// Assert
	assert.Equal(t, expected, result)
}
