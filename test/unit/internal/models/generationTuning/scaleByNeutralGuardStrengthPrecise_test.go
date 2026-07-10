package generationTuning_test

import (
	"math"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenValueIsScaledPrecisely_RoundsToThreeDecimalPlaces(t *testing.T) {
	// Arrange
	value := gofakeit.Float64Range(0.001, 10000)
	multiplier := gofakeit.Float64Range(0.1, 3)
	tuning := models.GenerationTuning{NeutralStackStrengthMultiplier: multiplier}
	expected := math.Round(value*multiplier*1000) / 1000

	// Act
	actual := tuning.ScaleByNeutralGuardStrengthPrecise(value)

	// Assert
	assert.InDelta(t, expected, actual, test_helpers.Delta)
}
