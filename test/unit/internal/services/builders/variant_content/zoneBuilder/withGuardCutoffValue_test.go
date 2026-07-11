package zoneBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenGuardCutoffValueIsProvided_SetsGuardCutoffValueOnBuiltZone(t *testing.T) {
	t.Parallel()
	// Arrange
	expectedCutoff := gofakeit.Number(1, 60000)
	builder := variant_content.NewZoneBuilder()

	// Act
	zone := builder.WithGuardCutoffValue(expectedCutoff).Build()

	// Assert
	assert.Equal(t, entities.Zone{GuardCutoffValue: expectedCutoff}, zone)
}
