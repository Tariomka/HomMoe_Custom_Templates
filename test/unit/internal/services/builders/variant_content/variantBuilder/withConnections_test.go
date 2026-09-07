package variantBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenConnectionsAreProvidedTwice_AppendsAllConnectionsOnBuiltVariant(t *testing.T) {
	t.Parallel()
	// Arrange
	firstConnection := template_model.Connection{Name: gofakeit.Word()}
	secondConnection := template_model.Connection{Name: gofakeit.Word()}
	thirdConnection := template_model.Connection{Name: gofakeit.Word()}
	builder := variant_content.NewVariantBuilder()

	// Act
	variant := builder.
		WithConnections(firstConnection, secondConnection).
		WithConnections(thirdConnection).
		Build()

	// Assert
	assert.Equal(t, template_model.Variant{
		Connections: []template_model.Connection{firstConnection, secondConnection, thirdConnection},
	}, variant)
}
