package guiHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
)

func TestWhenZoneIsNeutral_ReturnsServiceEquivalentDecision(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	zoneName := "Neutral-C"
	playerZoneNames := map[string]bool{"Spawn-A": true}
	expected := test_helpers.NewZoneEditorService().CanDeleteZone(zoneName, playerZoneNames)

	// Act
	result := handler.CanDeleteZone(zoneName, playerZoneNames)

	// Assert
	assert.Equal(t, expected, result)
}
