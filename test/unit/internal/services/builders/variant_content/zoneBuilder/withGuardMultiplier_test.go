package zoneBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenGuardMultiplierIsProvided_SetsGuardMultiplierOnBuiltZone(t *testing.T) {
	t.Parallel()
	// Arrange
	expectedMultiplier := gofakeit.Float64Range(0.1, 5)
	builder := variant_content.NewZoneBuilder()

	// Act
	zone := builder.WithGuardMultiplier(expectedMultiplier).Build()

	// Assert
	assert.Equal(t, entities.Zone{GuardMultiplier: expectedMultiplier}, zone)
}
