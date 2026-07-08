package placementRuleBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/placement_rule"
	"github.com/stretchr/testify/assert"
)

func TestWhenCrossroadsRuleIsBuilt_SetsCrossroadsTypeDistanceAndWeight(t *testing.T) {
	// Arrange
	builder := placement_rule.NewPlacementRuleBuilder()

	// Act
	rule := builder.BuildCrossroadsRule(placement_rule.Distance{Min: 0.25, Max: 0.5}, 2)

	// Assert
	assert.Equal(t, entities.PlacementRule{
		Type:      "Crossroads",
		TargetMin: 0.25,
		TargetMax: 0.5,
		Weight:    2,
	}, rule)
}
