package mainObjectBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenPlacementArgumentsAreProvidedTwice_AppendsAllPlacementArgumentsOnBuiltObject(t *testing.T) {
	// Arrange
	firstArgument := gofakeit.Word()
	secondArgument := gofakeit.Word()
	thirdArgument := gofakeit.Word()
	builder := variant_content.NewObjectBuilder()

	// Act
	mainObject := builder.
		WithPlacementArgs(firstArgument, secondArgument).
		WithPlacementArgs(thirdArgument).
		Build()

	// Assert
	assert.Equal(t, entities.MainObject{
		PlacementArgs: []string{firstArgument, secondArgument, thirdArgument},
	}, mainObject)
}
