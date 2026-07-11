package vec2_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenMultipliedByScalar_ScalesBothComponents(t *testing.T) {
	t.Parallel()
	// Arrange
	xComponent := gofakeit.Float64Range(-100, 100)
	yComponent := gofakeit.Float64Range(-100, 100)
	scalar := gofakeit.Float64Range(-10, 10)
	vector := data.NewVec2(xComponent, yComponent)

	// Act
	scaled := vector.MultiplyScalar(scalar)

	// Assert
	assert.Equal(t, data.Vec2[float64]{X: xComponent * scalar, Y: yComponent * scalar}, scaled)
}

func TestWhenMultipliedByZeroScalar_ReturnsZeroVector(t *testing.T) {
	t.Parallel()
	// Arrange
	vector := data.NewVec2(gofakeit.Float64Range(-100, 100), gofakeit.Float64Range(-100, 100))

	// Act
	scaled := vector.MultiplyScalar(0)

	// Assert
	assert.Equal(t, data.Vec2[float64]{X: 0, Y: 0}, scaled)
}
