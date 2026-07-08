package generationTuning_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenValueIsPositive_ScalesByNeutralStackStrengthMultiplier(t *testing.T) {
	// Arrange
	value := gofakeit.Number(1, 100000)
	multiplier := gofakeit.Float64Range(0.1, 3)
	tuning := models.GenerationTuning{NeutralStackStrengthMultiplier: multiplier}
	expected := int(float64(value) * multiplier)

	// Act
	scaled := tuning.ScaleByNeutralGuardStrength(value)

	// Assert
	assert.Equal(t, expected, scaled)
}

func TestWhenScaledNeutralGuardValueIsNegative_ClampsToZero(t *testing.T) {
	// Arrange
	tuning := models.GenerationTuning{NeutralStackStrengthMultiplier: 1.0}

	// Act
	scaled := tuning.ScaleByNeutralGuardStrength(-gofakeit.Number(1, 100000))

	// Assert
	assert.Equal(t, 0, scaled)
}
