package generationTuning_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenValueIsPositive_ScalesByResourceDensityMultiplier(t *testing.T) {
	// Arrange
	value := gofakeit.Float64Range(1, 10000)
	multiplier := gofakeit.Float64Range(0.1, 3)
	tuning := models.GenerationTuning{ResourceDensityMultiplier: multiplier}
	expected := int(value * multiplier)

	// Act
	scaled := tuning.ScaleByResourceDensity(value)

	// Assert
	assert.Equal(t, expected, scaled)
}

func TestWhenScaledResourceValueIsNegative_ClampsToZero(t *testing.T) {
	// Arrange
	tuning := models.GenerationTuning{ResourceDensityMultiplier: 1.0}

	// Act
	scaled := tuning.ScaleByResourceDensity(-gofakeit.Float64Range(1, 10000))

	// Assert
	assert.Equal(t, 0, scaled)
}
