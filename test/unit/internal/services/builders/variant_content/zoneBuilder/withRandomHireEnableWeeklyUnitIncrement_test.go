package zoneBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenRandomHireWeeklyIncrementFlagsAreProvided_SetsRandomHireWeeklyIncrementFlagsOnBuiltZone(t *testing.T) {
	// Arrange
	expectedFlags := []bool{gofakeit.Bool(), gofakeit.Bool(), gofakeit.Bool()}
	builder := variant_content.NewZoneBuilder()

	// Act
	zone := builder.WithRandomHireEnableWeeklyUnitIncrement(expectedFlags).Build()

	// Assert
	assert.Equal(t, entities.Zone{RandomHireEnableWeeklyUnitIncrement: expectedFlags}, zone)
}
