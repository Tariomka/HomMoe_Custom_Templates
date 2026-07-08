package ruleGuarded_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/content_rules"
	"github.com/stretchr/testify/assert"
)

func TestWhenDescriptionIsQueried_ExplainsGuardOverride(t *testing.T) {
	// Arrange
	rule := content_rules.NewRuleGuarded(true)

	// Act
	description := rule.Description()

	// Assert
	assert.Equal(t,
		"Forces the content item to be guarded or unguarded, regardless of the default behavior.",
		description)
}
