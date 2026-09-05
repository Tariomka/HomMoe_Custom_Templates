package connectionBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/stretchr/testify/assert"
)

func TestWhenCenterGatePlacementIsChosen_SetsCenterGatePlacementOnBuiltConnection(t *testing.T) {
	t.Parallel()
	// Arrange
	builder := variant_content.NewConnectionBuilder()

	// Act
	connection := builder.WithGatePlacementCenter().Build()

	// Assert
	assert.Equal(t, template_model.Connection{GatePlacement: "Center"}, connection)
}
