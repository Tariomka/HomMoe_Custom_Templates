package connectionBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenGuardZoneIsProvided_SetsGuardZoneOnBuiltConnection(t *testing.T) {
	// Arrange
	expectedGuardZone := gofakeit.Word()
	builder := variant_content.NewConnectionBuilder()

	// Act
	connection := builder.WithGuardZone(expectedGuardZone).Build()

	// Assert
	assert.Equal(t, entities.Connection{GuardZone: expectedGuardZone}, connection)
}
