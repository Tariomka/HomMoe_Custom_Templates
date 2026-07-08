package zoneBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenNameIsProvided_SetsNameOnBuiltZone(t *testing.T) {
	// Arrange
	expectedName := gofakeit.Word()
	builder := variant_content.NewZoneBuilder()

	// Act
	zone := builder.WithName(expectedName).Build()

	// Assert
	assert.Equal(t, entities.Zone{Name: expectedName}, zone)
}
