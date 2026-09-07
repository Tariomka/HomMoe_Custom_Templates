package typedRefBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenArgumentsAreProvidedTwice_AppendsAllArgumentsOnBuiltReference(t *testing.T) {
	t.Parallel()
	// Arrange
	firstArgument := gofakeit.Word()
	secondArgument := gofakeit.Word()
	thirdArgument := gofakeit.Word()
	builder := variant_content.NewRefBuilder()

	// Act
	reference := builder.WithArgs(firstArgument, secondArgument).WithArgs(thirdArgument).Build()

	// Assert
	assert.Equal(t, template_model.TypedRef{
		Args: []string{firstArgument, secondArgument, thirdArgument},
	}, reference)
}
