package guiHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
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
	expected := test_helpers.NewZoneEditorService().CountZoneCastles(zone)

	// Act
	result := handler.CountZoneCastles(zone)

	// Assert
	assert.Equal(t, expected, result)
}
