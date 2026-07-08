package zoneBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenUnguardedContentValueIsProvided_SetsUnguardedContentValueOnBuiltZone(t *testing.T) {
	// Arrange
	expectedValue := gofakeit.Number(1, 100000)
	builder := variant_content.NewZoneBuilder()

	// Act
	zone := builder.WithUnguardedContentValue(expectedValue).Build()

	// Assert
	assert.Equal(t, entities.Zone{UnguardedContentValue: expectedValue}, zone)
}
