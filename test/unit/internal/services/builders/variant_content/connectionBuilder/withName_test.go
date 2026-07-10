package connectionBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenNameIsProvided_SetsNameOnBuiltConnection(t *testing.T) {
	t.Parallel()
	// Arrange
	expectedName := gofakeit.Word()
	builder := variant_content.NewConnectionBuilder()

	// Act
	connection := builder.WithName(expectedName).Build()

	// Assert
	assert.Equal(t, entities.Connection{Name: expectedName}, connection)
}
