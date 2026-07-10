package ruleDistanceToRoad_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/content_rules"
	"github.com/stretchr/testify/assert"
)

func TestWhenDistanceIsNil_DefaultsToMedium(t *testing.T) {
	t.Parallel()
	// Arrange

	// Act
	rule := content_rules.NewRuleDistanceToRoad(nil)

	// Assert
	assert.Equal(t, content_rules.DistanceMedium, rule.Distance)
}

func TestWhenDistanceIsSupplied_UsesIt(t *testing.T) {
	t.Parallel()
	// Arrange
	distance := content_rules.DistanceFar

	// Act
	rule := content_rules.NewRuleDistanceToRoad(&distance)

	// Assert
	assert.Equal(t, content_rules.DistanceFar, rule.Distance)
}
