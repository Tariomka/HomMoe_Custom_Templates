package vec2_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenDividedByScalar_DividesBothComponents(t *testing.T) {
	t.Parallel()
	// Arrange
	xComponent := gofakeit.Float64Range(-100, 100)
	yComponent := gofakeit.Float64Range(-100, 100)
	scalar := gofakeit.Float64Range(1, 10)
	vector := data.NewVec2(xComponent, yComponent)

	// Act
	quotient := vector.DivideScalar(scalar)

	// Assert
	assert.Equal(t, data.Vec2[float64]{X: xComponent / scalar, Y: yComponent / scalar}, quotient)
}

func TestWhenIntVectorDividedByScalar_TruncatesQuotients(t *testing.T) {
	t.Parallel()
	// Arrange
	vector := data.NewVec2(9, 5)

	// Act
	quotient := vector.DivideScalar(2)

	// Assert
	assert.Equal(t, data.Vec2[int]{X: 4, Y: 2}, quotient)
}
