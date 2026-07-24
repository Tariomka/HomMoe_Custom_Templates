package ruleDistanceToRoad_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/content_rules"
	"github.com/stretchr/testify/assert"
)

func TestWhenDistanceIsFar_ShowsRuleNameWithDistanceName(t *testing.T) {
	t.Parallel()
	// Arrange
	distance := content_rules.DistanceVariation{Name: "Far", Min: 0.5, Max: 0.75}
	rule := content_rules.NewRuleDistanceToRoad(&distance)

	// Act
	displayText := rule.DisplayText()

	// Assert
	assert.Equal(t, "Distance to road: Far", displayText)
}
