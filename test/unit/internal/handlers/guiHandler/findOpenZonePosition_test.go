package guiHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/connection_editor"
	"github.com/stretchr/testify/assert"
)

func TestWhenPositionsAreOccupied_ReturnsServiceEquivalentPosition(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	occupied := [][2]float64{{0.2, 0.2}, {0.5, 0.5}}
	expected := connection_editor.NewDefaultZoneEditorService().FindOpenPosition(occupied)

	// Act
	result := handler.FindOpenZonePosition(occupied)

	// Assert
	assert.Equal(t, expected, result)
}
