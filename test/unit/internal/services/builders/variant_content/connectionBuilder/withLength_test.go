package connectionBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenLengthIsProvided_SetsLengthOnBuiltConnection(t *testing.T) {
	t.Parallel()
	// Arrange
	expectedLength := gofakeit.Float64Range(0.1, 3)
	builder := variant_content.NewConnectionBuilder()

	// Act
	connection := builder.WithLength(expectedLength).Build()

	// Assert
	assert.Equal(t, template_model.Connection{Length: expectedLength}, connection)
}
