package closestAcrossComponents_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/geometry"
	"github.com/stretchr/testify/assert"
)

func TestWhenOnlyOneComponentExists_ReportsNoPair(t *testing.T) {
	t.Parallel()
	// Arrange
	positions := []data.Vec2[float64]{data.NewVec2(0.1, 0.1), data.NewVec2(0.9, 0.9)}

	// Act
	_, found := geometry.FindClosestAcrossComponents(positions, [][]int{{0, 1}})

	// Assert
	assert.False(t, found)
}

func TestWhenComponentListIsEmpty_ReturnsNegativeIndexes(t *testing.T) {
	t.Parallel()
	// Arrange
	positions := []data.Vec2[float64]{data.NewVec2(0.1, 0.1)}

	// Act
	indexes, _ := geometry.FindClosestAcrossComponents(positions, nil)

	// Assert
	assert.Equal(t, data.NewVec2(-1, -1), indexes)
}

func TestWhenComponentsAreDisconnected_ReturnsClosestCrossComponentPair(t *testing.T) {
	t.Parallel()
	// Arrange
	positions := []data.Vec2[float64]{
		data.NewVec2(0.0, 0.0),
		data.NewVec2(0.4, 0.4),
		data.NewVec2(0.5, 0.5),
		data.NewVec2(1.0, 1.0),
	}

	// Act
	indexes, _ := geometry.FindClosestAcrossComponents(positions, [][]int{{0, 1}, {2, 3}})

	// Assert
	assert.Equal(t, data.NewVec2(1, 2), indexes)
}
