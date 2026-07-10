package connectionBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenFromZoneIsProvided_SetsFromOnBuiltConnection(t *testing.T) {
	t.Parallel()
	// Arrange
	expectedFrom := gofakeit.Word()
	builder := variant_content.NewConnectionBuilder()

	// Act
	connection := builder.WithFrom(expectedFrom).Build()

	// Assert
	assert.Equal(t, entities.Connection{From: expectedFrom}, connection)
}
