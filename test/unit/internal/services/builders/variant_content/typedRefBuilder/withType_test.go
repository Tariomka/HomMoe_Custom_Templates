package typedRefBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenTypeIsProvided_SetsTypeOnBuiltReference(t *testing.T) {
	// Arrange
	expectedType := gofakeit.Word()
	builder := variant_content.NewRefBuilder()

	// Act
	reference := builder.WithType(expectedType).Build()

	// Assert
	assert.Equal(t, entities.TypedRef{Type: expectedType}, reference)
}
