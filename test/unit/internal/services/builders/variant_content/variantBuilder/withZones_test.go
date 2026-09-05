package variantBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenZonesAreProvidedTwice_AppendsAllZonesOnBuiltVariant(t *testing.T) {
	t.Parallel()
	// Arrange
	firstZone := template_model.Zone{Name: gofakeit.Word()}
	secondZone := template_model.Zone{Name: gofakeit.Word()}
	thirdZone := template_model.Zone{Name: gofakeit.Word()}
	builder := variant_content.NewVariantBuilder()

	// Act
	variant := builder.WithZones(firstZone, secondZone).WithZones(thirdZone).Build()

	// Assert
	assert.Equal(t, template_model.Variant{
		Zones: []template_model.Zone{firstZone, secondZone, thirdZone},
	}, variant)
}
