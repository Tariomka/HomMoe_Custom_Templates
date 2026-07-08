package typedRefBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenArgumentsAreProvidedTwice_AppendsAllArgumentsOnBuiltReference(t *testing.T) {
	// Arrange
	firstArgument := gofakeit.Word()
	secondArgument := gofakeit.Word()
	thirdArgument := gofakeit.Word()
	builder := variant_content.NewRefBuilder()

	// Act
	reference := builder.WithArgs(firstArgument, secondArgument).WithArgs(thirdArgument).Build()

	// Assert
	assert.Equal(t, entities.TypedRef{
		Args: []string{firstArgument, secondArgument, thirdArgument},
	}, reference)
}
