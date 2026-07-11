package zoneBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenMatchMainObjectBiomeIsChosen_SetsAllThreeBiomesToMatchMainObject(t *testing.T) {
	t.Parallel()
	// Arrange
	expectedArgument := gofakeit.Word()
	expectedBiome := entities.TypedRef{Type: "MatchMainObject", Args: []string{expectedArgument}}
	builder := variant_content.NewZoneBuilder()

	// Act
	zone := builder.WithBiomeMatchMainObject(expectedArgument).Build()

	// Assert
	assert.Equal(t, entities.Zone{
		ZoneBiome:        expectedBiome,
		ContentBiome:     expectedBiome,
		MetaObjectsBiome: expectedBiome,
	}, zone)
}
