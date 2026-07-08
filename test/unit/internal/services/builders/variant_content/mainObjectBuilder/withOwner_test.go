package mainObjectBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenOwnerIsProvided_SetsOwnerOnBuiltObject(t *testing.T) {
	// Arrange
	expectedOwner := gofakeit.Word()
	builder := variant_content.NewObjectBuilder()

	// Act
	mainObject := builder.WithOwner(expectedOwner).Build()

	// Assert
	assert.Equal(t, entities.MainObject{Owner: expectedOwner}, mainObject)
}
