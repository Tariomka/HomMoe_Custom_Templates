package vec2_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenComponentsAreFuzzed_ReturnsDifferenceOfCrossComponentProducts(t *testing.T) {
	// Arrange
	firstX := gofakeit.Float64Range(-100, 100)
	firstY := gofakeit.Float64Range(-100, 100)
	secondX := gofakeit.Float64Range(-100, 100)
	secondY := gofakeit.Float64Range(-100, 100)
	firstVector := data.NewVec2(firstX, firstY)
	secondVector := data.NewVec2(secondX, secondY)

	// Act
	actual := firstVector.CrossProduct(secondVector)

	// Assert
	assert.InDelta(t, firstX*secondY-firstY*secondX, actual, test_helpers.Delta)
}

func TestWhenVectorsAreParallel_ReturnsZero(t *testing.T) {
	// Arrange
	xComponent := gofakeit.Float64Range(1, 100)
	yComponent := gofakeit.Float64Range(1, 100)
	vector := data.NewVec2(xComponent, yComponent)
	parallelVector := data.NewVec2(xComponent*2, yComponent*2)

	// Act
	actual := vector.CrossProduct(parallelVector)

	// Assert
	assert.InDelta(t, 0.0, actual, test_helpers.Delta)
}
