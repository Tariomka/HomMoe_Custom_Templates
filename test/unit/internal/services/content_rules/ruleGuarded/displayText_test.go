package ruleGuarded_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/content_rules"
	"github.com/stretchr/testify/assert"
)

func TestWhenGuarded_ShowsTrueState(t *testing.T) {
	// Arrange
	rule := content_rules.NewRuleGuarded(true)

	// Act
	displayText := rule.DisplayText()

	// Assert
	assert.Equal(t, "Guarded: true", displayText)
}

func TestWhenExplicitlyUnguarded_ShowsFalseState(t *testing.T) {
	// Arrange
	rule := content_rules.NewRuleGuarded(false)

	// Act
	displayText := rule.DisplayText()

	// Assert
	assert.Equal(t, "Guarded: false", displayText)
}
