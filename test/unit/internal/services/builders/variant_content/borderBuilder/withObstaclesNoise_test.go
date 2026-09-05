package borderBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenObstaclesNoiseIsProvided_SetsSingleObstaclesNoiseOnBuiltBorder(t *testing.T) {
	t.Parallel()
	// Arrange
	expectedAmplitude := gofakeit.Float64Range(0.01, 5)
	expectedFrequency := gofakeit.Number(1, 20)
	builder := variant_content.NewBorderBuilder()

	// Act
	border := builder.WithObstaclesNoise(expectedAmplitude, expectedFrequency).Build()

	// Assert
	assert.Equal(t, template_model.Border{
		ObstaclesNoise: []template_model.Noise{{Amplitude: expectedAmplitude, Frequency: expectedFrequency}},
	}, border)
}
