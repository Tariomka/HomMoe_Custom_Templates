package adjacency_test

import (
	"slices"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/stretchr/testify/assert"
)

func TestWhenGraphHasTwoComponents_ReturnsBothComponentsWithTheirMembers(t *testing.T) {
	t.Parallel()
	// Arrange
	nodes := []int{0, 1, 2, 3, 4}
	adjacency := data.NewAdjacency(nodes)
	adjacency.Link(0, 1)
	adjacency.Link(1, 2)
	adjacency.Link(3, 4)

	// Act
	components := adjacency.ConnectedComponents(nodes)

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
	adjacency := data.NewAdjacency(nodes)

	// Act
	components := adjacency.ConnectedComponents(nodes)

	// Assert
	assert.Equal(t, [][]int{{0}, {1}, {2}}, components)
}
