package variantBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenZonesAreProvidedTwice_AppendsAllZonesOnBuiltVariant(t *testing.T) {
	t.Parallel()
	// Arrange
	firstZone := entities.Zone{Name: gofakeit.Word()}
	secondZone := entities.Zone{Name: gofakeit.Word()}
	thirdZone := entities.Zone{Name: gofakeit.Word()}
	builder := variant_content.NewVariantBuilder()

	// Act
	variant := builder.WithZones(firstZone, secondZone).WithZones(thirdZone).Build()

	// Assert
	assert.Equal(t, entities.Variant{
		Zones: []entities.Zone{firstZone, secondZone, thirdZone},
	}, variant)
}
