package connectionBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenLengthIsProvided_SetsLengthOnBuiltConnection(t *testing.T) {
	// Arrange
	expectedLength := gofakeit.Float64Range(0.1, 3)
	builder := variant_content.NewConnectionBuilder()

	// Act
	connection := builder.WithLength(expectedLength).Build()

	// Assert
	assert.Equal(t, entities.Connection{Length: expectedLength}, connection)
}
