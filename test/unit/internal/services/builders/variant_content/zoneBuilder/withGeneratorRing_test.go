package zoneBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenGeneratorRingIsProvided_SetsGeneratorRingPointerOnBuiltZone(t *testing.T) {
	// Arrange
	expectedRing := gofakeit.Number(0, 5)
	builder := variant_content.NewZoneBuilder()

	// Act
	zone := builder.WithGeneratorRing(expectedRing).Build()

	// Assert
	assert.Equal(t, entities.Zone{GeneratorRing: &expectedRing}, zone)
}
