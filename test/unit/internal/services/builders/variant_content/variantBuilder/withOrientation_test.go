package variantBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenOrientationIsProvided_SetsOrientationOnBuiltVariant(t *testing.T) {
	t.Parallel()
	// Arrange
	expectedOrientation := template_model.Orientation{
		Mode:          gofakeit.Word(),
		ZeroAngleZone: gofakeit.Word(),
		BaseAngleMin:  gofakeit.Number(1, 179),
		BaseAngleMax:  gofakeit.Number(180, 359),
	}
	builder := variant_content.NewVariantBuilder()

	// Act
	variant := builder.WithOrientation(expectedOrientation).Build()

	// Assert
	assert.Equal(t, template_model.Variant{Orientation: expectedOrientation}, variant)
}
