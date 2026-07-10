package connectionBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenGuardMatchGroupIsProvided_SetsGuardMatchGroupOnBuiltConnection(t *testing.T) {
	t.Parallel()
	// Arrange
	expectedGroup := gofakeit.Word()
	builder := variant_content.NewConnectionBuilder()

	// Act
	connection := builder.WithGuardMatchGroup(expectedGroup).Build()

	// Assert
	assert.Equal(t, entities.Connection{GuardMatchGroup: expectedGroup}, connection)
}
