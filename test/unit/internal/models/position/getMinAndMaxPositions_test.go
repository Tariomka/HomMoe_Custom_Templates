package position_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestWhenPositionsAreScattered_ReturnsComponentWiseMinAndMax(t *testing.T) {
	// Arrange - min/max come from different positions per axis.
	positions := models.Positions{
		data.NewVec2(0.3, 0.8),
		data.NewVec2(0.1, 0.9),
		data.NewVec2(0.7, 0.2),
	}
	expected := []models.Position{data.NewVec2(0.1, 0.2), data.NewVec2(0.7, 0.9)}

	// Act
	minPosition, maxPosition := positions.GetMinAndMaxPositions()

	// Assert
	assert.Equal(t, expected, []models.Position{minPosition, maxPosition})
}
