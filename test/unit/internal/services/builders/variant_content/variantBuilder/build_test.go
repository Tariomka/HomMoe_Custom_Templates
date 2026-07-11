package variantBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenMultipleOptionsAreChained_ReturnsVariantWithAllAccumulatedValues(t *testing.T) {
	t.Parallel()
	// Arrange
	expectedOrientation := entities.Orientation{Mode: gofakeit.Word()}
	expectedBorder := entities.Border{CornerRadius: gofakeit.Float64Range(0.01, 1)}
	expectedZone := entities.Zone{Name: gofakeit.Word()}
	expectedConnection := entities.Connection{Name: gofakeit.Word()}
	builder := variant_content.NewVariantBuilder()

	// Act
	variant := builder.
		WithOrientation(expectedOrientation).
		WithBorder(expectedBorder).
		WithZones(expectedZone).
		WithConnections(expectedConnection).
		Build()

	// Assert
	assert.Equal(t, entities.Variant{
		Orientation: expectedOrientation,
		Border:      expectedBorder,
		Zones:       []entities.Zone{expectedZone},
		Connections: []entities.Connection{expectedConnection},
	}, variant)
}
