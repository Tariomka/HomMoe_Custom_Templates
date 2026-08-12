package guiHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
)

func TestWhenPositionsAreOccupied_ReturnsServiceEquivalentPosition(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	occupied := [][2]float64{{0.2, 0.2}, {0.5, 0.5}}
	expected := test_helpers.NewZoneEditorService().FindOpenPosition(occupied)

	// Act
	result := handler.FindOpenZonePosition(occupied)

	// Assert
	assert.Equal(t, expected, result)
}
