package position_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestWhenOnlyOneComponentExists_ReportsNoPair(t *testing.T) {
	t.Parallel()
	// Arrange
	positions := models.Positions{data.NewVec2(0.1, 0.1), data.NewVec2(0.9, 0.9)}

	// Act
	_, found := positions.GetShortestDistanceIndex([][]int{{0, 1}})

	// Assert
	assert.False(t, found)
}

func TestWhenComponentListIsEmpty_ReturnsNegativeIndexes(t *testing.T) {
	t.Parallel()
	// Arrange
	positions := models.Positions{data.NewVec2(0.1, 0.1)}

	// Act
	indexes, _ := positions.GetShortestDistanceIndex(nil)

	// Assert
	assert.Equal(t, data.NewVec2(-1, -1), indexes)
}

func TestWhenTwoComponentsExist_PicksClosestCrossComponentPair(t *testing.T) {
	t.Parallel()
	// Arrange
	positions := models.Positions{
		data.NewVec2(0.0, 0.0),
		data.NewVec2(0.4, 0.4),
		data.NewVec2(0.5, 0.5),
		data.NewVec2(1.0, 1.0),
	}

	// Act
	indexes, _ := positions.GetShortestDistanceIndex([][]int{{0, 1}, {2, 3}})

	// Assert
	assert.Equal(t, data.NewVec2(1, 2), indexes)
}
