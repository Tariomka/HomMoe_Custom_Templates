package connectionBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/stretchr/testify/assert"
)

func TestWhenCenterGatePlacementIsChosen_SetsCenterGatePlacementOnBuiltConnection(t *testing.T) {
	// Arrange
	builder := variant_content.NewConnectionBuilder()

	// Act
	connection := builder.WithGatePlacementCenter().Build()

	// Assert
	assert.Equal(t, entities.Connection{GatePlacement: "Center"}, connection)
}
