package connection_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model/template_variant_model"
	"github.com/stretchr/testify/assert"
)

func TestWhenAConnectionIsFlattened_TheSerializableRoundTripIsPreserved(t *testing.T) {
	t.Parallel()
	// Arrange
	expected := newPopulatedConnection()
	expected.IsUserAdded = false
	model := newPopulatedConnection()

	// Act
	roundTrippedModel := template_variant_model.ToConnectionModel(template_variant_model.ToConnectionEntity(model))

	// Assert
	assert.Equal(t, expected, roundTrippedModel)
}
