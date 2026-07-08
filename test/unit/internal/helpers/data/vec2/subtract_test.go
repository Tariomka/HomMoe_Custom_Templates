package vec2_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenVectorSubtracted_ReturnsComponentWiseDifference(t *testing.T) {
	// Arrange
	firstX := gofakeit.Float64Range(-1000, 1000)
	firstY := gofakeit.Float64Range(-1000, 1000)
	secondX := gofakeit.Float64Range(-1000, 1000)
	secondY := gofakeit.Float64Range(-1000, 1000)
	firstVector := data.NewVec2(firstX, firstY)
	secondVector := data.NewVec2(secondX, secondY)

	// Act
	difference := firstVector.Subtract(secondVector)

	// Assert
	assert.Equal(t, data.Vec2[float64]{X: firstX - secondX, Y: firstY - secondY}, difference)
}

func TestWhenVectorSubtractedFromItself_ReturnsZeroVector(t *testing.T) {
	// Arrange
	vector := data.NewVec2(gofakeit.Float64Range(-1000, 1000), gofakeit.Float64Range(-1000, 1000))

	// Act
	difference := vector.Subtract(vector)

	// Assert
	assert.Equal(t, data.Vec2[float64]{X: 0, Y: 0}, difference)
}
