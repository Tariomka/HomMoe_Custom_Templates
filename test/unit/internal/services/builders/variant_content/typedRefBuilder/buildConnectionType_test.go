package typedRefBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenConnectionReferenceIsBuilt_SetsConnectionTypeWithArguments(t *testing.T) {
	// Arrange
	expectedArgument := gofakeit.Word()
	builder := variant_content.NewRefBuilder()

	// Act
	reference := builder.BuildConnectionType(expectedArgument)

	// Assert
	assert.Equal(t, entities.TypedRef{Type: "Connection", Args: []string{expectedArgument}}, reference)
}
