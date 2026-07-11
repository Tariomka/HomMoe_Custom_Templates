package zoneBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenGeneratorPositionIsProvided_SetsGeneratorPositionPointerOnBuiltZone(t *testing.T) {
	t.Parallel()
	// Arrange
	expectedPosition := [2]float64{gofakeit.Float64Range(0, 1), gofakeit.Float64Range(0, 1)}
	builder := variant_content.NewZoneBuilder()

	// Act
	zone := builder.WithGeneratorPosition(expectedPosition).Build()

	// Assert
	assert.Equal(t, entities.Zone{GeneratorPosition: &expectedPosition}, zone)
}
