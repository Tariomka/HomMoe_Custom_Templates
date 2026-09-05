package connectionBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenGuardZoneIsProvided_SetsGuardZoneOnBuiltConnection(t *testing.T) {
	t.Parallel()
	// Arrange
	expectedGuardZone := gofakeit.Word()
	builder := variant_content.NewConnectionBuilder()

	// Act
	connection := builder.WithGuardZone(expectedGuardZone).Build()

	// Assert
	assert.Equal(t, template_model.Connection{GuardZone: expectedGuardZone}, connection)
}
