package variantBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenOrientationIsProvided_SetsOrientationOnBuiltVariant(t *testing.T) {
	// Arrange
	expectedOrientation := entities.Orientation{
		Mode:          gofakeit.Word(),
		ZeroAngleZone: gofakeit.Word(),
		BaseAngleMin:  gofakeit.Number(1, 179),
		BaseAngleMax:  gofakeit.Number(180, 359),
	}
	builder := variant_content.NewVariantBuilder()

	// Act
	variant := builder.WithOrientation(expectedOrientation).Build()

	// Assert
	assert.Equal(t, entities.Variant{Orientation: expectedOrientation}, variant)
}
