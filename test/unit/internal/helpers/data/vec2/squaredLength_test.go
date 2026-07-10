package vec2_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenComponentsAreFuzzed_ReturnsSumOfSquaredComponents(t *testing.T) {
	// Arrange
	xComponent := gofakeit.Float64Range(-100, 100)
	yComponent := gofakeit.Float64Range(-100, 100)
	vector := data.NewVec2(xComponent, yComponent)

	// Act
	actual := vector.SquaredLength()

	// Assert
	assert.InDelta(t, xComponent*xComponent+yComponent*yComponent, actual, test_helpers.Delta)
}

func TestWhenVectorIsZero_ReturnsZero(t *testing.T) {
	// Arrange
	vector := data.NewVec2(0.0, 0.0)

	// Act
	actual := vector.SquaredLength()

	// Assert
	assert.InDelta(t, 0.0, actual, test_helpers.Delta)
}
