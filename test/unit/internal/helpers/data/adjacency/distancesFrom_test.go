package adjacency_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/stretchr/testify/assert"
)

func TestWhenGraphIsChain_ComputesBreadthFirstDistances(t *testing.T) {
	t.Parallel()
	// Arrange
	adjacency := data.Adjacency[string]{}
	adjacency.Link("A", "B")
	adjacency.Link("B", "C")
	adjacency.Link("C", "D")
	expected := map[string]int{"A": 0, "B": 1, "C": 2, "D": 3}

	// Act
	distances := adjacency.DistancesFrom("A")

	// Assert
	assert.Equal(t, expected, distances)
}

func TestWhenGraphIsDisconnected_OmitsUnreachableNodes(t *testing.T) {
	t.Parallel()
	// Arrange
	adjacency := data.Adjacency[string]{}
	adjacency.Link("A", "B")
	adjacency.Link("C", "D")

	// Act
	distances := adjacency.DistancesFrom("A")

	// Assert
	assert.Equal(t, map[string]int{"A": 0, "B": 1}, distances)
}

func TestWhenShorterPathExists_PrefersShortestDistance(t *testing.T) {
	t.Parallel()
	// Arrange
	adjacency := data.Adjacency[string]{}
	adjacency.Link("A", "B")
	adjacency.Link("B", "C")
	adjacency.Link("A", "C")

	// Act
	distances := adjacency.DistancesFrom("A")

	// Assert
	assert.Equal(t, 1, distances["C"])
}
