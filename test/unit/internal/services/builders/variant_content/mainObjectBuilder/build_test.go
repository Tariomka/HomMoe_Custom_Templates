package mainObjectBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenMultipleOptionsAreChained_ReturnsObjectWithAllAccumulatedValues(t *testing.T) {
	t.Parallel()
	// Arrange
	expectedOwner := gofakeit.Word()
	expectedGuardValue := gofakeit.Number(1, 60000)
	builder := variant_content.NewObjectBuilder()

	// Act
	mainObject := builder.
		WithTypeCity().
		WithOwner(expectedOwner).
		WithGuardValue(expectedGuardValue).
		WithCastleQualityPoor().
		WithPlacementCenter().
		Build()

	// Assert
	assert.Equal(t, template_model.MainObject{
		Type:                     "City",
		Owner:                    expectedOwner,
		GuardValue:               expectedGuardValue,
		BuildingsConstructionSid: "poor_buildings_construction",
		Placement:                "Center",
	}, mainObject)
}
