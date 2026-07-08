package zoneBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenMultipleOptionsAreChained_ReturnsZoneWithAllAccumulatedValues(t *testing.T) {
	// Arrange
	expectedName := gofakeit.Word()
	expectedSize := gofakeit.Float64Range(0.1, 10)
	expectedGuardCutoff := gofakeit.Number(1, 60000)
	expectedPool := []string{gofakeit.Word(), gofakeit.Word()}
	builder := variant_content.NewZoneBuilder()

	// Act
	zone := builder.
		WithName(expectedName).
		WithSize(expectedSize).
		WithLayoutSpawns().
		WithGuardCutoffValue(expectedGuardCutoff).
		WithGuardedContentPool(expectedPool).
		Build()

	// Assert
	assert.Equal(t, entities.Zone{
		Name:               expectedName,
		Size:               expectedSize,
		Layout:             "zone_layout_spawns",
		GuardCutoffValue:   expectedGuardCutoff,
		GuardedContentPool: expectedPool,
	}, zone)
}
