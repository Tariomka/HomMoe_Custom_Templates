package placementRuleBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/placement_rule"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenRoadRuleIsBuilt_SetsRoadTypeDistanceAndWeight(t *testing.T) {
	t.Parallel()
	// Arrange
	expectedDistance := placement_rule.Distance{
		Min: gofakeit.Float64Range(0.01, 0.4),
		Max: gofakeit.Float64Range(0.5, 0.95),
	}
	expectedWeight := gofakeit.Number(1, 100)
	builder := placement_rule.NewPlacementRuleBuilder()

	// Act
	rule := builder.BuildRoadRule(expectedDistance, expectedWeight)

	// Assert
	assert.Equal(t, entities.PlacementRule{
		Type:      "Road",
		TargetMin: expectedDistance.Min,
		TargetMax: expectedDistance.Max,
		Weight:    expectedWeight,
	}, rule)
}
