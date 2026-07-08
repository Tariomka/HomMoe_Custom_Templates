package placementRuleBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/placement_rule"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenCrossroadsRuleIsBuilt_SetsCrossroadsTypeDistanceAndWeight(t *testing.T) {
	// Arrange
	expectedDistance := placement_rule.Distance{
		Min: gofakeit.Float64Range(0.01, 0.4),
		Max: gofakeit.Float64Range(0.5, 0.95),
	}
	expectedWeight := gofakeit.Number(1, 100)
	builder := placement_rule.NewPlacementRuleBuilder()

	// Act
	rule := builder.BuildCrossroadsRule(expectedDistance, expectedWeight)

	// Assert
	assert.Equal(t, entities.PlacementRule{
		Type:      "Crossroads",
		TargetMin: expectedDistance.Min,
		TargetMax: expectedDistance.Max,
		Weight:    expectedWeight,
	}, rule)
}
