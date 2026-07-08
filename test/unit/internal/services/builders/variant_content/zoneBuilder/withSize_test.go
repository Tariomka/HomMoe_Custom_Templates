package zoneBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenSizeIsProvided_SetsSizeOnBuiltZone(t *testing.T) {
	// Arrange
	expectedSize := gofakeit.Float64Range(0.1, 10)
	builder := variant_content.NewZoneBuilder()

	// Act
	zone := builder.WithSize(expectedSize).Build()

	// Assert
	assert.Equal(t, entities.Zone{Size: expectedSize}, zone)
}
