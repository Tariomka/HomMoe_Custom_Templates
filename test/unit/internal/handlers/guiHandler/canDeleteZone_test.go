package guiHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/handlers"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/connection_editor"
	"github.com/stretchr/testify/assert"
)

func TestWhenZoneIsNeutral_ReturnsServiceEquivalentDecision(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := handlers.NewGuiHandler()
	zoneName := "Neutral-C"
	playerZoneNames := map[string]bool{"Spawn-A": true}
	expected := connection_editor.NewZoneEditorService().CanDeleteZone(zoneName, playerZoneNames)

	// Act
	result := handler.CanDeleteZone(zoneName, playerZoneNames)

	// Assert
	assert.Equal(t, expected, result)
}
