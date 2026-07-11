package typedRefBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenBiomeMatchMainObjectReferenceIsBuilt_SetsMatchMainObjectTypeWithArguments(t *testing.T) {
	t.Parallel()
	// Arrange
	expectedArgument := gofakeit.Word()
	builder := variant_content.NewRefBuilder()

	// Act
	reference := builder.BuildBiomeMatchMainObjectType(expectedArgument)

	// Assert
	assert.Equal(t, entities.TypedRef{Type: "MatchMainObject", Args: []string{expectedArgument}}, reference)
}
