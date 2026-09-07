package variantBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenBorderIsProvided_SetsBorderOnBuiltVariant(t *testing.T) {
	t.Parallel()
	// Arrange
	expectedBorder := template_model.Border{
		CornerRadius:   gofakeit.Float64Range(0.01, 1),
		ObstaclesWidth: gofakeit.Number(1, 20),
		WaterWidth:     gofakeit.Number(1, 20),
	}
	builder := variant_content.NewVariantBuilder()

	// Act
	variant := builder.WithBorder(expectedBorder).Build()

	// Assert
	assert.Equal(t, template_model.Variant{Border: expectedBorder}, variant)
}
