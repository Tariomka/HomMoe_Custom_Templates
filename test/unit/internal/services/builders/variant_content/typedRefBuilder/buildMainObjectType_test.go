package typedRefBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenMainObjectReferenceIsBuilt_SetsMainObjectTypeWithArguments(t *testing.T) {
	t.Parallel()
	// Arrange
	expectedArgument := gofakeit.Word()
	builder := variant_content.NewRefBuilder()

	// Act
	reference := builder.BuildMainObjectType(expectedArgument)

	// Assert
	assert.Equal(t, template_model.TypedRef{Type: "MainObject", Args: []string{expectedArgument}}, reference)
}
