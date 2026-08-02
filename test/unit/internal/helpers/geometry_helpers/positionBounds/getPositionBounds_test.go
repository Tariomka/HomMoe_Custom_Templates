package positionBounds_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/geometry_helpers"
	"github.com/stretchr/testify/assert"
)

func TestWhenPositionsAreScattered_ReturnsComponentWiseBounds(t *testing.T) {
	t.Parallel()
	// Arrange
	positions := []data.Vec2[float64]{
		data.NewVec2(0.3, 0.8),
		data.NewVec2(0.1, 0.9),
		data.NewVec2(0.7, 0.2),
	}
	expected := []data.Vec2[float64]{data.NewVec2(0.1, 0.2), data.NewVec2(0.7, 0.9)}

	// Act
	minimumPosition, maximumPosition := geometry_helpers.GetPositionBounds(positions)

	// Assert
	assert.Equal(t, expected, []data.Vec2[float64]{minimumPosition, maximumPosition})
}
