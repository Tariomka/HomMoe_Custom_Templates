package vec2_test

import (
	"image"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/stretchr/testify/assert"
)

func TestWhenFractionIsAboveHalf_RoundsComponentsUp(t *testing.T) {
	t.Parallel()
	// Arrange
	vector := data.NewVec2(4.6, 2.7)

	// Act
	point := vector.ToPointRounded()

	// Assert
	assert.Equal(t, image.Pt(5, 3), point)
}

func TestWhenFractionIsBelowHalf_RoundsComponentsDown(t *testing.T) {
	t.Parallel()
	// Arrange
	vector := data.NewVec2(4.4, 2.2)

	// Act
	point := vector.ToPointRounded()

	// Assert
	assert.Equal(t, image.Pt(4, 2), point)
}

func TestWhenFractionIsExactlyHalf_RoundsComponentsAwayFromZero(t *testing.T) {
	t.Parallel()
	// Arrange
	vector := data.NewVec2(4.5, -2.5)

	// Act
	point := vector.ToPointRounded()

	// Assert
	assert.Equal(t, image.Pt(5, -3), point)
}
