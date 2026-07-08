package zoneEditor_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/connection_editor"
	"github.com/stretchr/testify/assert"
)

func TestWhenBoardIsEmpty_PicksCornerFarthestFromCenter(t *testing.T) {
	// Arrange

	// Act
	position := connection_editor.FindOpenPosition(nil)

	// Assert
	assert.InDeltaSlice(t, []float64{0.9, 0.9}, position[:], 0.0001,
		"with no occupied positions a grid corner maximizes the distance to the center; float rounding favors the last corner")
}

func TestWhenCornerIsCrowded_PicksPositionAwayFromIt(t *testing.T) {
	// Arrange
	occupied := [][2]float64{{0.1, 0.1}, {0.1, 0.2}, {0.2, 0.1}}

	// Act
	position := connection_editor.FindOpenPosition(occupied)

	// Assert
	assert.Greater(t, position[0]+position[1], 1.0,
		"the best spot should be far from the cluttered top-left corner")
}
