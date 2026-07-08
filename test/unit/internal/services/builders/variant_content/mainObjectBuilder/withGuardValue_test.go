package mainObjectBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenGuardValueIsProvided_SetsGuardValueOnBuiltObject(t *testing.T) {
	// Arrange
	expectedGuardValue := gofakeit.Number(1, 60000)
	builder := variant_content.NewObjectBuilder()

	// Act
	mainObject := builder.WithGuardValue(expectedGuardValue).Build()

	// Assert
	assert.Equal(t, entities.MainObject{GuardValue: expectedGuardValue}, mainObject)
}
