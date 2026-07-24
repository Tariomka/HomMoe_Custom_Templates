package adjacency_test

import (
	"slices"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/graph"
	"github.com/stretchr/testify/assert"
)

func TestWhenGraphHasTwoComponents_ReturnsBothComponentsWithTheirMembers(t *testing.T) {
	t.Parallel()
	// Arrange
	nodes := []int{0, 1, 2, 3, 4}
	adjacency := graph.NewAdjacency(nodes)
	graph.Link(adjacency, 0, 1)
	graph.Link(adjacency, 1, 2)
	graph.Link(adjacency, 3, 4)

	// Act
	components := graph.ConnectedComponents(adjacency, nodes)

	// Assert
	sortedComponents := make([][]int, 0, len(components))
	for _, component := range components {
		sorted := append([]int(nil), component...)
		slices.Sort(sorted)
		sortedComponents = append(sortedComponents, sorted)
	}
	assert.Equal(t, [][]int{{0, 1, 2}, {3, 4}}, sortedComponents)
}

func TestWhenNoLinksExist_ReturnsOneComponentPerNode(t *testing.T) {
	t.Parallel()
	// Arrange
	nodes := []int{0, 1, 2}
	adjacency := graph.NewAdjacency(nodes)

	// Act
	components := graph.ConnectedComponents(adjacency, nodes)

	// Assert
	assert.Equal(t, [][]int{{0}, {1}, {2}}, components)
}
