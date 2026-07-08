package vec2_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenVectorsMultipliedComponentWise_ReturnsComponentWiseProduct(t *testing.T) {
	// Arrange
	firstX := gofakeit.Float64Range(-100, 100)
	firstY := gofakeit.Float64Range(-100, 100)
	secondX := gofakeit.Float64Range(-100, 100)
	secondY := gofakeit.Float64Range(-100, 100)
	firstVector := data.NewVec2(firstX, firstY)
	secondVector := data.NewVec2(secondX, secondY)

	// Act
	product := firstVector.MultiplyComponent(secondVector)

	// Assert
	assert.Equal(t, data.Vec2[float64]{X: firstX * secondX, Y: firstY * secondY}, product)
}

func TestWhenMultipliedByOnesVector_ReturnsReceiverUnchanged(t *testing.T) {
	// Arrange
	xComponent := gofakeit.Float64Range(-100, 100)
	yComponent := gofakeit.Float64Range(-100, 100)
	vector := data.NewVec2(xComponent, yComponent)
	onesVector := data.NewVec2(1.0, 1.0)

	// Act
	product := vector.MultiplyComponent(onesVector)

	// Assert
	assert.Equal(t, data.Vec2[float64]{X: xComponent, Y: yComponent}, product)
}
