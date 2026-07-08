package placementRuleBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/placement_rule"
	"github.com/stretchr/testify/assert"
)

func TestWhenCastleRuleIsBuilt_TargetsPrimaryMainObjectWithDistanceAndWeight(t *testing.T) {
	// Arrange
	builder := placement_rule.NewPlacementRuleBuilder()

	// Act
	rule := builder.BuildCastleRule(placement_rule.Distance{Min: 0.1, Max: 0.25}, 1)

	// Assert
	assert.Equal(t, entities.PlacementRule{
		Type:      "MainObject",
		Args:      []any{"0"},
		TargetMin: 0.1,
		TargetMax: 0.25,
		Weight:    1,
	}, rule)
}
