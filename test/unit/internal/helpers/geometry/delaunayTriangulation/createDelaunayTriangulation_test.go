package delaunayTriangulation_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/geometry"
	"github.com/stretchr/testify/assert"
)

func TestWhenNoPositionsAreProvided_ReturnsNoEdges(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	edges := geometry.CreateDelaunayTriangulation(nil)

	// Assert
	assert.Nil(t, edges)
}

func TestWhenOnePositionIsProvided_ReturnsNoEdges(t *testing.T) {
	t.Parallel()
	// Arrange
	positions := []data.Vec2[float64]{data.NewVec2(0.5, 0.5)}

	// Act
	edges := geometry.CreateDelaunayTriangulation(positions)

	// Assert
	assert.Nil(t, edges)
}

func TestWhenTwoPositionsAreProvided_ReturnsSingleEdge(t *testing.T) {
	t.Parallel()
	// Arrange
	positions := []data.Vec2[float64]{data.NewVec2(0.1, 0.1), data.NewVec2(0.9, 0.9)}

	// Act
	edges := geometry.CreateDelaunayTriangulation(positions)

	// Assert
	assert.Equal(t, []data.Vec2[int]{data.NewVec2(0, 1)}, edges)
}

func TestWhenThreePositionsFormTriangle_ReturnsNormalizedEdgesInOrder(t *testing.T) {
	t.Parallel()
	// Arrange
	positions := []data.Vec2[float64]{
		data.NewVec2(0.1, 0.1),
		data.NewVec2(0.9, 0.1),
		data.NewVec2(0.5, 0.9),
	}
	expected := []data.Vec2[int]{data.NewVec2(0, 1), data.NewVec2(0, 2), data.NewVec2(1, 2)}

	// Act
	edges := geometry.CreateDelaunayTriangulation(positions)

	// Assert
	assert.Equal(t, expected, edges)
}

func TestWhenTriangleWindingIsReversed_ReturnsAllTriangleEdges(t *testing.T) {
	t.Parallel()
	// Arrange
	positions := []data.Vec2[float64]{
		data.NewVec2(0.1, 0.1),
		data.NewVec2(0.5, 0.9),
		data.NewVec2(0.9, 0.1),
	}
	expected := []data.Vec2[int]{data.NewVec2(0, 1), data.NewVec2(0, 2), data.NewVec2(1, 2)}

	// Act
	edges := geometry.CreateDelaunayTriangulation(positions)

	// Assert
	assert.Equal(t, expected, edges)
}

func TestWhenFourPositionsFormSquare_ReturnsDeterministicEdges(t *testing.T) {
	t.Parallel()
	// Arrange
	positions := []data.Vec2[float64]{
		data.NewVec2(0.0, 0.0),
		data.NewVec2(1.0, 0.0),
		data.NewVec2(1.0, 1.0),
		data.NewVec2(0.0, 1.0),
	}
	expected := geometry.CreateDelaunayTriangulation(positions)

	// Act
	actual := geometry.CreateDelaunayTriangulation(positions)

	// Assert
	assert.Equal(t, expected, actual)
}

func TestWhenPositionsAreCollinear_ReturnsNoTriangleEdges(t *testing.T) {
	t.Parallel()
	// Arrange
	positions := []data.Vec2[float64]{
		data.NewVec2(0.0, 0.0),
		data.NewVec2(0.5, 0.5),
		data.NewVec2(1.0, 1.0),
	}

	// Act
	edges := geometry.CreateDelaunayTriangulation(positions)

	// Assert
	assert.Empty(t, edges)
}

func TestWhenPositionsContainDuplicate_ReturnsNoInvalidIndexes(t *testing.T) {
	t.Parallel()
	// Arrange
	positions := []data.Vec2[float64]{
		data.NewVec2(0.0, 0.0),
		data.NewVec2(0.0, 0.0),
		data.NewVec2(1.0, 0.0),
		data.NewVec2(0.0, 1.0),
	}

	// Act
	edges := geometry.CreateDelaunayTriangulation(positions)

	// Assert
	invalidCount := 0
	for _, edge := range edges {
		if edge.X < 0 || edge.Y >= len(positions) || edge.X >= edge.Y {
			invalidCount++
		}
	}
	assert.Zero(t, invalidCount)
}
