package zoneAdjacency_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenIndexAdjacencyIsCreated_AllocatesEmptyNeighbourSetPerIndex(t *testing.T) {
	t.Parallel()
	// Arrange
	size := gofakeit.Number(1, 10)
	expected := map[int]map[int]bool{}
	for index := range size {
		expected[index] = map[int]bool{}
	}

	// Act
	adjacency := models.NewZoneIndexAdjacency(size)

	// Assert
	assert.Equal(t, models.ZoneIndexAdjacency(expected), *adjacency)
}
