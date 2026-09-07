package borderBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenWaterWidthIsProvided_SetsWaterWidthOnBuiltBorder(t *testing.T) {
	t.Parallel()
	// Arrange
	expectedWidth := gofakeit.Number(1, 20)
	builder := variant_content.NewBorderBuilder()

	// Act
	border := builder.WithWaterWidth(expectedWidth).Build()

	// Assert
	assert.Equal(t, template_model.Border{WaterWidth: expectedWidth}, border)
}
