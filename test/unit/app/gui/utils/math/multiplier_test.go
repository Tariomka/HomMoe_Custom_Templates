package math_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/utils"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenFactorAndBaseAreApplied_LinearCombinationIsReturned(t *testing.T) {
	// Arrange
	value := gofakeit.Float32Range(0, 1)
	base := gofakeit.Float32Range(-10, 10)
	factor := gofakeit.Float32Range(-10, 10)

	// Act
	result := utils.Multiplier(value, base, factor)

	// Assert
	assert.InDelta(t, value*factor+base, result, test_helpers.Delta)
}

func TestWhenFactorIsZero_BaseIsReturned(t *testing.T) {
	// Arrange
	base := gofakeit.Float32Range(-10, 10)

	// Act
	result := utils.Multiplier(gofakeit.Float32Range(0, 1), base, 0)

	// Assert
	assert.InDelta(t, base, result, test_helpers.Delta)
}
