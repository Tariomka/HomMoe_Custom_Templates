package ruleDistanceToRoad_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/content_rules"
	"github.com/stretchr/testify/assert"
)

func TestWhenDistanceIsFar_ShowsRuleNameWithDistanceName(t *testing.T) {
	// Arrange
	distance := content_rules.DistanceFar
	rule := content_rules.NewRuleDistanceToRoad(&distance)

	// Act
	displayText := rule.DisplayText()

	// Assert
	assert.Equal(t, "Distance to road: Far", displayText)
}
