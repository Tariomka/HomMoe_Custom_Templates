package ruleDistanceToRoad_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/content_rules"
	"github.com/stretchr/testify/assert"
)

func TestWhenNameIsQueried_ReturnsDistanceToRoad(t *testing.T) {
	// Arrange
	rule := content_rules.NewRuleDistanceToRoad(nil)

	// Act
	name := rule.Name()

	// Assert
	assert.Equal(t, "Distance to road", name)
}
