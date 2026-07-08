package zoneBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenLayoutNameIsProvided_SetsLayoutOnBuiltZone(t *testing.T) {
	// Arrange
	expectedLayout := gofakeit.Word()
	builder := variant_content.NewZoneBuilder()

	// Act
	zone := builder.WithLayout(expectedLayout).Build()

	// Assert
	assert.Equal(t, entities.Zone{Layout: expectedLayout}, zone)
}
