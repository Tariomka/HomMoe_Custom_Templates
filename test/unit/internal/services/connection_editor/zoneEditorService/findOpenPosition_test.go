package zoneEditorService_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
)

func TestWhenBoardIsEmpty_PicksCornerFarthestFromCenter(t *testing.T) {
	t.Parallel()
	// Arrange

	// Act
	position := test_helpers.NewZoneEditorService().FindOpenPosition(nil)

	// Assert
	assert.InDeltaSlice(
		t,
		[]float64{0.9, 0.9},
		[]float64{position.X, position.Y},
		0.0001,
		"with no occupied positions a grid corner maximizes the distance to the center; float rounding favors the last corner",
	)
}

func TestWhenCornerIsCrowded_PicksPositionAwayFromIt(t *testing.T) {
	t.Parallel()
	// Arrange
	occupied := []data.Vec2[float64]{{X: 0.1, Y: 0.1}, {X: 0.1, Y: 0.2}, {X: 0.2, Y: 0.1}}

	// Act
	position := test_helpers.NewZoneEditorService().FindOpenPosition(occupied)

	// Assert
	assert.Greater(t, position.X+position.Y, 1.0,
		"the best spot should be far from the cluttered top-left corner")
}
