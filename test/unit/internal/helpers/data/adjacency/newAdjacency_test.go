package adjacency_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/stretchr/testify/assert"
)

func TestWhenNodesAreProvided_AllocatesEmptyNeighbourSetPerNode(t *testing.T) {
	t.Parallel()
	// Arrange
	nodes := []int{0, 1, 2}
	expected := data.Adjacency[int]{0: {}, 1: {}, 2: {}}

	// Act
	adjacency := data.NewAdjacency(nodes)

	// Assert
	assert.Equal(t, expected, adjacency)
}
