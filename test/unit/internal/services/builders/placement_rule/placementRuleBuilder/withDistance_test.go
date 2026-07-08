package placementRuleBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/placement_rule"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenDistanceIsProvided_SetsTargetRangeOnBuiltRule(t *testing.T) {
	// Arrange
	expectedDistance := placement_rule.Distance{
		Min: gofakeit.Float64Range(0.01, 0.4),
		Max: gofakeit.Float64Range(0.5, 0.95),
	}
	builder := placement_rule.NewPlacementRuleBuilder()

	// Act
	rule := builder.WithDistance(expectedDistance).Build()

	// Assert
	assert.Equal(t, entities.PlacementRule{
		TargetMin: expectedDistance.Min,
		TargetMax: expectedDistance.Max,
	}, rule)
}
