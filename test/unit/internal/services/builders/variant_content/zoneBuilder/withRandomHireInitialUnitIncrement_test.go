package zoneBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenRandomHireInitialIncrementsAreProvided_SetsRandomHireInitialIncrementsOnBuiltZone(t *testing.T) {
	t.Parallel()
	// Arrange
	expectedIncrements := []int{gofakeit.Number(1, 100), gofakeit.Number(1, 100)}
	builder := variant_content.NewZoneBuilder()

	// Act
	zone := builder.WithRandomHireInitialUnitIncrement(expectedIncrements).Build()

	// Assert
	assert.Equal(t, template_model.Zone{RandomHireInitialUnitIncrement: expectedIncrements}, zone)
}
