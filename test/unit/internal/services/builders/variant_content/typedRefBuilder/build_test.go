package typedRefBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenTypeAndArgumentsAreChained_ReturnsReferenceWithAllAccumulatedValues(t *testing.T) {
	// Arrange
	expectedType := gofakeit.Word()
	expectedArgument := gofakeit.Word()
	builder := variant_content.NewRefBuilder()

	// Act
	reference := builder.WithType(expectedType).WithArgs(expectedArgument).Build()

	// Assert
	assert.Equal(t, entities.TypedRef{Type: expectedType, Args: []string{expectedArgument}}, reference)
}
