package zoneBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenRandomHireWeeklyIncrementFlagsAreProvided_SetsRandomHireWeeklyIncrementFlagsOnBuiltZone(t *testing.T) {
	t.Parallel()
	// Arrange
	expectedFlags := []bool{gofakeit.Bool(), gofakeit.Bool(), gofakeit.Bool()}
	builder := variant_content.NewZoneBuilder()

	// Act
	zone := builder.WithRandomHireEnableWeeklyUnitIncrement(expectedFlags).Build()

	// Assert
	assert.Equal(t, template_model.Zone{RandomHireEnableWeeklyUnitIncrement: expectedFlags}, zone)
}
