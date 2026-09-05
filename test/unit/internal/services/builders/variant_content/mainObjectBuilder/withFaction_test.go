package mainObjectBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenFactionTypeAndArgumentsAreProvided_SetsFactionReferenceOnBuiltObject(t *testing.T) {
	t.Parallel()
	// Arrange
	expectedType := gofakeit.Word()
	firstArgument := gofakeit.Word()
	secondArgument := gofakeit.Word()
	builder := variant_content.NewObjectBuilder()

	// Act
	mainObject := builder.WithFaction(expectedType, firstArgument, secondArgument).Build()

	// Assert
	assert.Equal(t, template_model.MainObject{
		Faction: &template_model.TypedRef{Type: expectedType, Args: []string{firstArgument, secondArgument}},
	}, mainObject)
}
