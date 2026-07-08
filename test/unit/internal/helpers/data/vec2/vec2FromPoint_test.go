package vec2_test

import (
	"image"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenPointGiven_ReturnsFloatVectorWithConvertedComponents(t *testing.T) {
	// Arrange
	xCoordinate := gofakeit.Number(-1000, 1000)
	yCoordinate := gofakeit.Number(-1000, 1000)
	point := image.Pt(xCoordinate, yCoordinate)

	// Act
	vector := data.Vec2FromPoint[float64](point)

	// Assert
	assert.Equal(t, data.Vec2[float64]{X: float64(xCoordinate), Y: float64(yCoordinate)}, vector)
}

func TestWhenPointGiven_ReturnsIntVectorWithSameComponents(t *testing.T) {
	// Arrange
	xCoordinate := gofakeit.Number(-1000, 1000)
	yCoordinate := gofakeit.Number(-1000, 1000)
	point := image.Pt(xCoordinate, yCoordinate)

	// Act
	vector := data.Vec2FromPoint[int](point)

	// Assert
	assert.Equal(t, data.Vec2[int]{X: xCoordinate, Y: yCoordinate}, vector)
}
