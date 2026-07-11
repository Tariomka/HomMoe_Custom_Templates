package zoneAdjacency_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestWhenTwoIndexesAreLinked_CreatesSymmetricIndexAdjacency(t *testing.T) {
	t.Parallel()
	// Arrange
	adjacency := models.NewZoneIndexAdjacency(3)
	expected := models.ZoneIndexAdjacency{
		0: {2: true},
		1: {},
		2: {0: true},
	}

	// Act
	adjacency.Link(0, 2)

	// Assert
	assert.Equal(t, expected, *adjacency)
}
