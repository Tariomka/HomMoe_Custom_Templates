package placementRuleBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/placement_rule"
	"github.com/stretchr/testify/assert"
)

func TestWhenRoadRuleIsBuilt_SetsRoadTypeDistanceAndWeight(t *testing.T) {
	// Arrange
	builder := placement_rule.NewPlacementRuleBuilder()

	// Act
	rule := builder.BuildRoadRule(placement_rule.Distance{Min: 0.1, Max: 0.25}, 1)

	// Assert
	assert.Equal(t, entities.PlacementRule{
		Type:      "Road",
		TargetMin: 0.1,
		TargetMax: 0.25,
		Weight:    1,
	}, rule)
}
