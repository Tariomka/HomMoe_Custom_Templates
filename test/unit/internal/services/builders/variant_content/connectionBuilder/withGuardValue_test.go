package connectionBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenGuardValueIsProvided_SetsGuardValueOnBuiltConnection(t *testing.T) {
	t.Parallel()
	// Arrange
	expectedGuardValue := gofakeit.Number(1, 60000)
	builder := variant_content.NewConnectionBuilder()

	// Act
	connection := builder.WithGuardValue(expectedGuardValue).Build()

	// Assert
	assert.Equal(t, template_model.Connection{GuardValue: expectedGuardValue}, connection)
}
