package mainObjectBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenInitialUnitIncrementIsProvided_SetsInitialUnitIncrementOnBuiltObject(t *testing.T) {
	t.Parallel()
	// Arrange
	expectedIncrement := gofakeit.Number(1, 100)
	builder := variant_content.NewObjectBuilder()

	// Act
	mainObject := builder.WithInitialUnitIncrement(expectedIncrement).Build()

	// Assert
	assert.Equal(t, entities.MainObject{InitialUnitIncrement: expectedIncrement}, mainObject)
}
