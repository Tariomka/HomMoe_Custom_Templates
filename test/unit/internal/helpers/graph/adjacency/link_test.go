package adjacency_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/graph"
	"github.com/stretchr/testify/assert"
)

func TestWhenTwoNodesAreLinked_CreatesSymmetricAdjacency(t *testing.T) {
	t.Parallel()
	// Arrange
	adjacency := graph.Adjacency[string]{}
	expected := graph.Adjacency[string]{"A": {"B": true}, "B": {"A": true}}

	// Act
	graph.Link(adjacency, "A", "B")

	// Assert
	assert.Equal(t, expected, adjacency)
}

func TestWhenNodeIsLinkedTwice_KeepsBothNeighbours(t *testing.T) {
	t.Parallel()
	// Arrange
	adjacency := graph.Adjacency[string]{}
	graph.Link(adjacency, "A", "B")
	expected := map[string]bool{"B": true, "C": true}

	// Act
	graph.Link(adjacency, "A", "C")

	// Assert
	assert.Equal(t, expected, adjacency["A"])
}
