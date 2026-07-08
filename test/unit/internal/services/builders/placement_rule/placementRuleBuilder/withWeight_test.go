package placementRuleBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/placement_rule"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenWeightIsProvided_SetsWeightOnBuiltRule(t *testing.T) {
	// Arrange
	expectedWeight := gofakeit.Number(1, 100)
	builder := placement_rule.NewPlacementRuleBuilder()

	// Act
	rule := builder.WithWeight(expectedWeight).Build()

	// Assert
	assert.Equal(t, entities.PlacementRule{Weight: expectedWeight}, rule)
}
