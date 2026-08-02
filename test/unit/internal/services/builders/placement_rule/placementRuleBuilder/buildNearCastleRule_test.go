package placementRuleBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/placement_rule"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenNearCastleRuleIsBuilt_TargetsPrimaryMainObjectAtNearDistance(t *testing.T) {
	t.Parallel()
	// Arrange
	expectedWeight := gofakeit.Number(1, 100)
	builder := placement_rule.NewPlacementRuleBuilder()

	// Act
	rule := builder.BuildNearCastleRule(expectedWeight)

	// Assert
	assert.Equal(t, entities.PlacementRule{
		Type:      "MainObject",
		Args:      []any{"0"},
		TargetMin: 0.1,
		TargetMax: 0.25,
		Weight:    expectedWeight,
	}, rule)
}
