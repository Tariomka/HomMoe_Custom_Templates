package mainObjectBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenSpawnIsProvided_SetsSpawnOnBuiltObject(t *testing.T) {
	t.Parallel()
	// Arrange
	expectedSpawn := gofakeit.Word()
	builder := variant_content.NewObjectBuilder()

	// Act
	mainObject := builder.WithSpawn(expectedSpawn).Build()

	// Assert
	assert.Equal(t, entities.MainObject{Spawn: expectedSpawn}, mainObject)
}
