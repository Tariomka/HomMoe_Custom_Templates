package position_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/stretchr/testify/assert"
)

func normalizedEdges(edges []models.ConnectionIndexes) []models.ConnectionIndexes {
	normalized := make([]models.ConnectionIndexes, 0, len(edges))
	for _, edge := range edges {
		if edge.X > edge.Y {
			edge.X, edge.Y = edge.Y, edge.X
		}
		normalized = append(normalized, edge)
	}
	return normalized
}

func TestWhenListHasOnePosition_ReturnsNoEdges(t *testing.T) {
	// Arrange
	positions := models.Positions{data.NewVec2(0.5, 0.5)}

	// Act
	edges := positions.CreateDelaunayTriangulation()

	// Assert
	assert.Nil(t, edges)
}

func TestWhenListHasTwoPositions_ReturnsSingleEdge(t *testing.T) {
	// Arrange
	positions := models.Positions{data.NewVec2(0.1, 0.1), data.NewVec2(0.9, 0.9)}

	// Act
	edges := positions.CreateDelaunayTriangulation()

	// Assert
	assert.Equal(t, []models.ConnectionIndexes{data.NewVec2(0, 1)}, edges)
}

func TestWhenListHasThreePositions_ReturnsAllTriangleEdges(t *testing.T) {
	// Arrange
	positions := models.Positions{
		data.NewVec2(0.1, 0.1),
		data.NewVec2(0.9, 0.1),
		data.NewVec2(0.5, 0.9),
	}
	expected := []models.ConnectionIndexes{
		data.NewVec2(0, 1),
		data.NewVec2(0, 2),
		data.NewVec2(1, 2),
	}

	// Act
	edges := positions.CreateDelaunayTriangulation()

	// Assert
	assert.ElementsMatch(t, expected, normalizedEdges(edges))
}

func TestWhenFourPositionsFormASquare_OmitsOneDiagonal(t *testing.T) {
	// Arrange - a unit square triangulates into 4 hull edges + exactly 1 diagonal.
	positions := models.Positions{
		data.NewVec2(0.0, 0.0),
		data.NewVec2(1.0, 0.0),
		data.NewVec2(1.0, 1.0),
		data.NewVec2(0.0, 1.0),
	}

	// Act
	edges := positions.CreateDelaunayTriangulation()

	// Assert
	assert.Len(t, edges, 5)
}
