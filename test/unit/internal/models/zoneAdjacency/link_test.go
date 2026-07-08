package zoneAdjacency_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestWhenTwoLabelsAreLinked_CreatesSymmetricAdjacency(t *testing.T) {
	// Arrange
	adjacency := models.ZoneAdjacency{}
	expected := models.ZoneAdjacency{
		"A": {"B": true},
		"B": {"A": true},
	}

	// Act
	adjacency.Link("A", "B")

	// Assert
	assert.Equal(t, expected, adjacency)
}

func TestWhenLabelIsLinkedTwice_KeepsBothNeighbours(t *testing.T) {
	// Arrange
	adjacency := models.ZoneAdjacency{}
	adjacency.Link("A", "B")
	expected := map[string]bool{"B": true, "C": true}

	// Act
	adjacency.Link("A", "C")

	// Assert
	assert.Equal(t, expected, adjacency["A"])
}
