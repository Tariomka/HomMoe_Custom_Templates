package mainObjectBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenFactionsAreProvidedTwice_AppendsAllFactionsOnBuiltObject(t *testing.T) {
	t.Parallel()
	// Arrange
	firstFaction := gofakeit.Word()
	secondFaction := gofakeit.Word()
	thirdFaction := gofakeit.Word()
	builder := variant_content.NewObjectBuilder()

	// Act
	mainObject := builder.WithFactions(firstFaction, secondFaction).WithFactions(thirdFaction).Build()

	// Assert
	assert.Equal(t, template_model.MainObject{
		Factions: []string{firstFaction, secondFaction, thirdFaction},
	}, mainObject)
}
