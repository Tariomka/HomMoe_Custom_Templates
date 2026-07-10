package connectionBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/stretchr/testify/assert"
)

func TestWhenPortalTypeIsChosen_SetsPortalConnectionTypeOnBuiltConnection(t *testing.T) {
	t.Parallel()
	// Arrange
	builder := variant_content.NewConnectionBuilder()

	// Act
	connection := builder.WithConnectionTypePortal().Build()

	// Assert
	assert.Equal(t, entities.Connection{ConnectionType: "Portal"}, connection)
}
