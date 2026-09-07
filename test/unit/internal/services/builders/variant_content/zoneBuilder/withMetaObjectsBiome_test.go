package zoneBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenMetaObjectsBiomeIsProvided_SetsOnlyMetaObjectsBiomeOnBuiltZone(t *testing.T) {
	t.Parallel()
	// Arrange
	expectedBiome := template_model.TypedRef{Type: gofakeit.Word(), Args: []string{gofakeit.Word()}}
	builder := variant_content.NewZoneBuilder()

	// Act
	zone := builder.WithMetaObjectsBiome(expectedBiome).Build()

	// Assert
	assert.Equal(t, template_model.Zone{MetaObjectsBiome: expectedBiome}, zone)
}
