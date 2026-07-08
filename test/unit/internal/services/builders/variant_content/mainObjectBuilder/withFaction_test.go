package mainObjectBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenFactionTypeAndArgumentsAreProvided_SetsFactionReferenceOnBuiltObject(t *testing.T) {
	// Arrange
	expectedType := gofakeit.Word()
	firstArgument := gofakeit.Word()
	secondArgument := gofakeit.Word()
	builder := variant_content.NewObjectBuilder()

	// Act
	mainObject := builder.WithFaction(expectedType, firstArgument, secondArgument).Build()

	// Assert
	assert.Equal(t, entities.MainObject{
		Faction: &entities.TypedRef{Type: expectedType, Args: []string{firstArgument, secondArgument}},
	}, mainObject)
}
