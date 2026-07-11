package zoneBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenRoadsAreProvided_SetsRoadsOnBuiltZone(t *testing.T) {
	t.Parallel()
	// Arrange
	expectedRoads := []entities.Road{
		{Type: "Stone", From: entities.TypedRef{Type: gofakeit.Word()}},
		{Type: "Dirt", To: entities.TypedRef{Type: gofakeit.Word()}},
	}
	builder := variant_content.NewZoneBuilder()

	// Act
	zone := builder.WithRoads(expectedRoads).Build()

	// Assert
	assert.Equal(t, entities.Zone{Roads: expectedRoads}, zone)
}
