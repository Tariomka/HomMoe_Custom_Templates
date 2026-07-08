package variantBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenConnectionsAreProvidedTwice_AppendsAllConnectionsOnBuiltVariant(t *testing.T) {
	// Arrange
	firstConnection := entities.Connection{Name: gofakeit.Word()}
	secondConnection := entities.Connection{Name: gofakeit.Word()}
	thirdConnection := entities.Connection{Name: gofakeit.Word()}
	builder := variant_content.NewVariantBuilder()

	// Act
	variant := builder.
		WithConnections(firstConnection, secondConnection).
		WithConnections(thirdConnection).
		Build()

	// Assert
	assert.Equal(t, entities.Variant{
		Connections: []entities.Connection{firstConnection, secondConnection, thirdConnection},
	}, variant)
}
