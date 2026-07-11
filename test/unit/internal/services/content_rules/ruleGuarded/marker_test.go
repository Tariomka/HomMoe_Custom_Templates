package ruleGuarded_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/content_rules"
	"github.com/stretchr/testify/assert"
)

func TestWhenGuarded_ReturnsG(t *testing.T) {
	t.Parallel()
	// Arrange
	rule := content_rules.NewRuleGuarded(true)

	// Act
	marker := rule.Marker()

	// Assert
	assert.Equal(t, "G", marker)
}

func TestWhenExplicitlyUnguarded_ReturnsNegatedG(t *testing.T) {
	t.Parallel()
	// Arrange
	rule := content_rules.NewRuleGuarded(false)

	// Act
	marker := rule.Marker()

	// Assert
	assert.Equal(t, "!G", marker)
}
