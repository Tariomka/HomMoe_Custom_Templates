package adjacency_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/stretchr/testify/assert"
)

func TestWhenTwoNodesAreLinked_CreatesSymmetricAdjacency(t *testing.T) {
	t.Parallel()
	// Arrange
	adjacency := data.Adjacency[string]{}
	expected := data.Adjacency[string]{"A": {"B": true}, "B": {"A": true}}

	// Act
	adjacency.Link("A", "B")

	// Assert
	assert.Equal(t, expected, adjacency)
}

func TestWhenNodeIsLinkedTwice_KeepsBothNeighbours(t *testing.T) {
	t.Parallel()
	// Arrange
	adjacency := data.Adjacency[string]{}
	adjacency.Link("A", "B")
	expected := map[string]bool{"B": true, "C": true}

	// Act
	adjacency.Link("A", "C")

	// Assert
	assert.Equal(t, expected, adjacency["A"])
}
