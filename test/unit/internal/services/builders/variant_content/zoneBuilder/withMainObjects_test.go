package zoneBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenMainObjectsAreProvided_SetsMainObjectsOnBuiltZone(t *testing.T) {
	t.Parallel()
	// Arrange
	expectedObjects := []entities.MainObject{
		{Type: "Spawn", Owner: gofakeit.Word()},
		{Type: "City", GuardValue: gofakeit.Number(1, 60000)},
	}
	builder := variant_content.NewZoneBuilder()

	// Act
	zone := builder.WithMainObjects(expectedObjects).Build()

	// Assert
	assert.Equal(t, entities.Zone{MainObjects: expectedObjects}, zone)
}
