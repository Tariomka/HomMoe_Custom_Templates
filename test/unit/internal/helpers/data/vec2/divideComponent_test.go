package vec2_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenVectorsDividedComponentWise_ReturnsComponentWiseQuotient(t *testing.T) {
	t.Parallel()
	// Arrange
	firstX := gofakeit.Float64Range(-100, 100)
	firstY := gofakeit.Float64Range(-100, 100)
	secondX := gofakeit.Float64Range(1, 100)
	secondY := gofakeit.Float64Range(1, 100)
	firstVector := data.NewVec2(firstX, firstY)
	secondVector := data.NewVec2(secondX, secondY)

	// Act
	quotient := firstVector.DivideComponent(secondVector)

	// Assert
	assert.Equal(t, data.Vec2[float64]{X: firstX / secondX, Y: firstY / secondY}, quotient)
}

func TestWhenIntVectorDividedComponentWise_TruncatesQuotients(t *testing.T) {
	t.Parallel()
	// Arrange
	firstVector := data.NewVec2(7, 9)
	secondVector := data.NewVec2(2, 4)

	// Act
	quotient := firstVector.DivideComponent(secondVector)

	// Assert
	assert.Equal(t, data.Vec2[int]{X: 3, Y: 2}, quotient)
}
