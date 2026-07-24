package ruleDistanceToTown_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/content_rules"
	"github.com/stretchr/testify/assert"
)

func TestWhenDistanceIsNear_ShowsRuleNameWithDistanceName(t *testing.T) {
	t.Parallel()
	// Arrange
	distance := content_rules.DistanceVariation{Name: "Near", Min: 0.1, Max: 0.25}
	rule := content_rules.NewRuleDistanceToTown(&distance)

	// Act
	displayText := rule.DisplayText()

	// Assert
	assert.Equal(t, "Distance to town: Near", displayText)
}
