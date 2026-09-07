package guiHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
)

func TestWhenPositionsAreOccupied_ReturnsServiceEquivalentPosition(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	occupied := []data.Vec2[float64]{{X: 0.2, Y: 0.2}, {X: 0.5, Y: 0.5}}
	expected := test_helpers.NewZoneEditorService().FindOpenPosition(occupied)

	// Act
	result := handler.FindOpenZonePosition(occupied)

	// Assert
	assert.Equal(t, expected, result)
}
