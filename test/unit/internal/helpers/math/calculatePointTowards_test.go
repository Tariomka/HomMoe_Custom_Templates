package math_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/stretchr/testify/assert"
)

func TestWhenTheTwoPointsAreCoincident_ReportsNoDirection(t *testing.T) {
	t.Parallel()
	// Arrange
	from := data.NewVec2(120.5, 44.25)

	// Act
	_, ok := helpers.CalculatePointTowards(from, from, 10)

	// Assert
	assert.False(t, ok)
}

func TestWhenTheTwoPointsAreCoincident_ReturnsTheZeroVector(t *testing.T) {
	t.Parallel()
	// Arrange
	from := data.NewVec2(120.5, 44.25)

	// Act
	moved, _ := helpers.CalculatePointTowards(from, from, 10)

	// Assert
	assert.Equal(t, data.NewVec2(0.0, 0.0), moved)
}

func TestWhenTheTargetIsFartherThanAPixel_MovesAlongTheDirectionByTheDistance(t *testing.T) {
	t.Parallel()
	// Arrange
	from := data.NewVec2(100.0, 100.0)
	toward := data.NewVec2(130.0, 140.0)

	// Act
	moved, _ := helpers.CalculatePointTowards(from, toward, 10)

	// Assert
	assert.InDeltaSlice(t, []float64{106.0, 108.0}, []float64{moved.X, moved.Y}, 1e-9)
}

func TestWhenTheDistanceIsFractional_TheResultKeepsItsFractionalPart(t *testing.T) {
	t.Parallel()
	// Arrange
	from := data.NewVec2(0.0, 0.0)
	toward := data.NewVec2(10.0, 0.0)

	// Act
	moved, _ := helpers.CalculatePointTowards(from, toward, 0.25)

	// Assert
	assert.InDelta(t, 0.25, moved.X, 1e-9)
}
