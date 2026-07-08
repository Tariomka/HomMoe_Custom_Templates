package ruleDistanceToRoad_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/content_rules"
	"github.com/stretchr/testify/assert"
)

func TestWhenMarkerIsQueried_ReturnsR(t *testing.T) {
	// Arrange
	rule := content_rules.NewRuleDistanceToRoad(nil)

	// Act
	marker := rule.Marker()

	// Assert
	assert.Equal(t, "R", marker)
}
