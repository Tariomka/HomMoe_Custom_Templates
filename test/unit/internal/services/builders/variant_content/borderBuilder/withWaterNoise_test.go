package borderBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenWaterNoiseIsProvided_SetsSingleWaterNoiseOnBuiltBorder(t *testing.T) {
	// Arrange
	expectedAmplitude := gofakeit.Float64Range(0.01, 5)
	expectedFrequency := gofakeit.Number(1, 20)
	builder := variant_content.NewBorderBuilder()

	// Act
	border := builder.WithWaterNoise(expectedAmplitude, expectedFrequency).Build()

	// Assert
	assert.Equal(t, entities.Border{
		WaterNoise: []entities.Noise{{Amplitude: expectedAmplitude, Frequency: expectedFrequency}},
	}, border)
}
