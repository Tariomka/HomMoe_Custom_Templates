package variantBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenMultipleOptionsAreChained_ReturnsVariantWithAllAccumulatedValues(t *testing.T) {
	t.Parallel()
	// Arrange
	expectedOrientation := template_model.Orientation{Mode: gofakeit.Word()}
	expectedBorder := template_model.Border{CornerRadius: gofakeit.Float64Range(0.01, 1)}
	expectedZone := template_model.Zone{Name: gofakeit.Word()}
	expectedConnection := template_model.Connection{Name: gofakeit.Word()}
	builder := variant_content.NewVariantBuilder()

	// Act
	variant := builder.
		WithOrientation(expectedOrientation).
		WithBorder(expectedBorder).
		WithZones(expectedZone).
		WithConnections(expectedConnection).
		Build()

	// Assert
	assert.Equal(t, template_model.Variant{
		Orientation: expectedOrientation,
		Border:      expectedBorder,
		Zones:       []template_model.Zone{expectedZone},
		Connections: []template_model.Connection{expectedConnection},
	}, variant)
}
