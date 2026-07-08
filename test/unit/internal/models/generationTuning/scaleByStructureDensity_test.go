package generationTuning_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenValueIsPositive_ScalesByStructureDensityMultiplier(t *testing.T) {
	// Arrange
	value := gofakeit.Float64Range(1, 10000)
	multiplier := gofakeit.Float64Range(0.1, 3)
	tuning := models.GenerationTuning{StructureDensityMultiplier: multiplier}
	expected := int(value * multiplier)

	// Act
	scaled := tuning.ScaleByStructureDensity(value)

	// Assert
	assert.Equal(t, expected, scaled)
}

func TestWhenScaledStructureValueIsNegative_ClampsToZero(t *testing.T) {
	// Arrange
	tuning := models.GenerationTuning{StructureDensityMultiplier: 1.0}

	// Act
	scaled := tuning.ScaleByStructureDensity(-gofakeit.Float64Range(1, 10000))

	// Assert
	assert.Equal(t, 0, scaled)
}
