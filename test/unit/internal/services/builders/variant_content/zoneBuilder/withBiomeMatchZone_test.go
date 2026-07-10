package zoneBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenMatchZoneBiomeIsChosen_SetsAllThreeBiomesToMatchZone(t *testing.T) {
	t.Parallel()
	// Arrange
	expectedArgument := gofakeit.Word()
	expectedBiome := entities.TypedRef{Type: "MatchZone", Args: []string{expectedArgument}}
	builder := variant_content.NewZoneBuilder()

	// Act
	zone := builder.WithBiomeMatchZone(expectedArgument).Build()

	// Assert
	assert.Equal(t, entities.Zone{
		ZoneBiome:        expectedBiome,
		ContentBiome:     expectedBiome,
		MetaObjectsBiome: expectedBiome,
	}, zone)
}
