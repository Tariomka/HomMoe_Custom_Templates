package connection_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model/template_variant_model"
	"github.com/stretchr/testify/assert"
)

func TestWhenAConnectionEntityIsLifted_AllSerializableFieldsArePreserved(t *testing.T) {
	t.Parallel()
	// Arrange
	expected := newPopulatedConnection()
	expected.IsUserAdded = false
	entity := template_variant_model.ToConnectionEntity(newPopulatedConnection())

	// Act
	model := template_variant_model.ToConnectionModel(entity)

	// Assert
	assert.Equal(t, expected, model)
}
