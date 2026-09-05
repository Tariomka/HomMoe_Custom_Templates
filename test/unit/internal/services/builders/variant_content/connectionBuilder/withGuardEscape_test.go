package connectionBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenGuardEscapeIsProvided_SetsGuardEscapeOnBuiltConnection(t *testing.T) {
	t.Parallel()
	// Arrange
	expectedGuardEscape := gofakeit.Bool()
	builder := variant_content.NewConnectionBuilder()

	// Act
	connection := builder.WithGuardEscape(expectedGuardEscape).Build()

	// Assert
	assert.Equal(t, template_model.Connection{GuardEscape: expectedGuardEscape}, connection)
}
