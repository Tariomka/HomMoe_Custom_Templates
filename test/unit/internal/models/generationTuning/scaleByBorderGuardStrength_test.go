package generationTuning_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenValueIsPositive_ScalesByBorderGuardStrengthMultiplier(t *testing.T) {
	// Arrange
	value := gofakeit.Number(1, 100000)
	multiplier := gofakeit.Float64Range(0.1, 3)
	tuning := models.GenerationTuning{BorderGuardStrengthMultiplier: multiplier}
	expected := int(float64(value) * multiplier)

	// Act
	scaled := tuning.ScaleByBorderGuardStrength(value)

	// Assert
	assert.Equal(t, expected, scaled)
}

func TestWhenScaledBorderGuardValueIsNegative_ClampsToZero(t *testing.T) {
	// Arrange
	tuning := models.GenerationTuning{BorderGuardStrengthMultiplier: 1.0}

	// Act
	scaled := tuning.ScaleByBorderGuardStrength(-gofakeit.Number(1, 100000))

	// Assert
	assert.Equal(t, 0, scaled)
}
