package connectionBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenToPortalRulesAreProvidedTwice_AppendsAllToRulesOnBuiltConnection(t *testing.T) {
	// Arrange
	firstRule := entities.PlacementRule{Type: "Road", Weight: gofakeit.Number(1, 100)}
	secondRule := entities.PlacementRule{Type: "Crossroads", Weight: gofakeit.Number(1, 100)}
	thirdRule := entities.PlacementRule{Type: "MainObject", Weight: gofakeit.Number(1, 100)}
	builder := variant_content.NewConnectionBuilder()

	// Act
	connection := builder.
		WithPortalPlacementRulesTo(firstRule, secondRule).
		WithPortalPlacementRulesTo(thirdRule).
		Build()

	// Assert
	assert.Equal(t, entities.Connection{
		PortalPlacementRulesTo: []entities.PlacementRule{firstRule, secondRule, thirdRule},
	}, connection)
}
