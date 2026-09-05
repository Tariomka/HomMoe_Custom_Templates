package zoneBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenGuardRandomizationIsProvided_SetsGuardRandomizationOnBuiltZone(t *testing.T) {
	t.Parallel()
	// Arrange
	expectedRandomization := gofakeit.Float64Range(0.01, 1)
	builder := variant_content.NewZoneBuilder()

	// Act
	zone := builder.WithGuardRandomization(expectedRandomization).Build()

	// Assert
	assert.Equal(t, template_model.Zone{GuardRandomization: expectedRandomization}, zone)
}
