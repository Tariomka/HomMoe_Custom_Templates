package connectionBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenToPortalRulesAreProvidedTwice_AppendsAllToRulesOnBuiltConnection(t *testing.T) {
	t.Parallel()
	// Arrange
	firstRule := template_model.PlacementRule{Type: "Road", Weight: gofakeit.Number(1, 100)}
	secondRule := template_model.PlacementRule{Type: "Crossroads", Weight: gofakeit.Number(1, 100)}
	thirdRule := template_model.PlacementRule{Type: "MainObject", Weight: gofakeit.Number(1, 100)}
	builder := variant_content.NewConnectionBuilder()

	// Act
	connection := builder.
		WithPortalPlacementRulesTo(firstRule, secondRule).
		WithPortalPlacementRulesTo(thirdRule).
		Build()

	// Assert
	assert.Equal(t, template_model.Connection{
		PortalPlacementRulesTo: []template_model.PlacementRule{firstRule, secondRule, thirdRule},
	}, connection)
}
