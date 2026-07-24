package adjacency_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/graph"
	"github.com/stretchr/testify/assert"
)

func TestWhenGraphIsChain_ComputesBreadthFirstDistances(t *testing.T) {
	t.Parallel()
	// Arrange
	adjacency := graph.Adjacency[string]{}
	graph.Link(adjacency, "A", "B")
	graph.Link(adjacency, "B", "C")
	graph.Link(adjacency, "C", "D")
	expected := map[string]int{"A": 0, "B": 1, "C": 2, "D": 3}

	// Act
	distances := graph.DistancesFrom(adjacency, "A")

	// Assert
	assert.Equal(t, expected, distances)
}

func TestWhenGraphIsDisconnected_OmitsUnreachableNodes(t *testing.T) {
	t.Parallel()
	// Arrange
	adjacency := graph.Adjacency[string]{}
	graph.Link(adjacency, "A", "B")
	graph.Link(adjacency, "C", "D")

	// Act
	distances := graph.DistancesFrom(adjacency, "A")

	// Assert
	assert.Equal(t, map[string]int{"A": 0, "B": 1}, distances)
}

func TestWhenShorterPathExists_PrefersShortestDistance(t *testing.T) {
	t.Parallel()
	// Arrange
	adjacency := graph.Adjacency[string]{}
	graph.Link(adjacency, "A", "B")
	graph.Link(adjacency, "B", "C")
	graph.Link(adjacency, "A", "C")

	// Act
	distances := graph.DistancesFrom(adjacency, "A")

	// Assert
	assert.Equal(t, 1, distances["C"])
}
