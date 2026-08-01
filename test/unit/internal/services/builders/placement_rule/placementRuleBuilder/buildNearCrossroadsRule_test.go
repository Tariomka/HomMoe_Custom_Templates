package placementRuleBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/placement_rule"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenNearCrossroadsRuleIsBuilt_TargetsPortalNearDistance(t *testing.T) {
	t.Parallel()
	// Arrange
	expectedWeight := gofakeit.Number(1, 100)
	builder := placement_rule.NewPlacementRuleBuilder()

	// Act
	rule := builder.BuildNearCrossroadsRule(expectedWeight)

	// Assert
	assert.Equal(t, entities.PlacementRule{
		Type:      "Crossroads",
		TargetMin: 0.075,
		TargetMax: 0.35,
		Weight:    expectedWeight,
	}, rule)
}
