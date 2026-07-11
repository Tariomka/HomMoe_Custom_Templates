package zoneAdjacency_test

import (
	"slices"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestWhenGraphHasTwoComponents_ReturnsBothComponentsWithTheirMembers(t *testing.T) {
	t.Parallel()
	// Arrange
	adjacency := models.NewZoneIndexAdjacency(5)
	adjacency.Link(0, 1)
	adjacency.Link(1, 2)
	adjacency.Link(3, 4)

	// Act
	components := adjacency.FindIndexes(5)

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
	adjacency := models.NewZoneIndexAdjacency(3)

	// Act
	components := adjacency.FindIndexes(3)

	// Assert
	assert.Equal(t, [][]int{{0}, {1}, {2}}, components)
}
