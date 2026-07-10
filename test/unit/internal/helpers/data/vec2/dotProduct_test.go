package vec2_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenComponentsAreFuzzed_ReturnsSumOfComponentProducts(t *testing.T) {
	// Arrange
	firstX := gofakeit.Float64Range(-100, 100)
	firstY := gofakeit.Float64Range(-100, 100)
	secondX := gofakeit.Float64Range(-100, 100)
	secondY := gofakeit.Float64Range(-100, 100)
	firstVector := data.NewVec2(firstX, firstY)
	secondVector := data.NewVec2(secondX, secondY)

	// Act
	actual := firstVector.DotProduct(secondVector)

	// Assert
	assert.InDelta(t, firstX*secondX+firstY*secondY, actual, test_helpers.Delta)
}

func TestWhenVectorsArePerpendicular_ReturnsZero(t *testing.T) {
	// Arrange
	component := gofakeit.Float64Range(1, 100)
	firstVector := data.NewVec2(component, 0.0)
	secondVector := data.NewVec2(0.0, component)

	// Act
	actual := firstVector.DotProduct(secondVector)

	// Assert
	assert.InDelta(t, 0.0, actual, test_helpers.Delta)
}
