package ruleDistanceToTown_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/content_rules"
	"github.com/stretchr/testify/assert"
)

func TestWhenRuleIsApplied_AppendsMainObjectPlacementRuleWithDistanceBounds(t *testing.T) {
	// Arrange
	distance := content_rules.DistanceNear
	rule := content_rules.NewRuleDistanceToTown(&distance)
	item := entities.MandatoryContentItem{SID: "x"}

	// Act
	rule.Apply(&item)

	// Assert
	assert.Equal(t, []entities.PlacementRule{
		{Type: "MainObject", Args: []any{"0"}, TargetMin: 0.1, TargetMax: 0.25, Weight: 1},
	}, item.Rules)
}
