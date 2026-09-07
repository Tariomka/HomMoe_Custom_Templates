package zoneBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenDiplomacyModifierIsProvided_SetsDiplomacyModifierOnBuiltZone(t *testing.T) {
	t.Parallel()
	// Arrange
	expectedModifier := gofakeit.Float64Range(0.01, 5)
	builder := variant_content.NewZoneBuilder()

	// Act
	zone := builder.WithDiplomacyModifier(expectedModifier).Build()

	// Assert
	assert.Equal(t, template_model.Zone{DiplomacyModifier: expectedModifier}, zone)
}
