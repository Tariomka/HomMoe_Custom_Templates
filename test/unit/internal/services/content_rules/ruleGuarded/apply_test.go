package ruleGuarded_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/content_rules"
	"github.com/stretchr/testify/assert"
)

func TestWhenGuardedRuleIsApplied_SetsItemGuarded(t *testing.T) {
	t.Parallel()
	// Arrange
	rule := content_rules.NewRuleGuarded(true)
	item := template_model.MandatoryContentItem{SID: "x"}

	// Act
	rule.Apply(&item)

	// Assert
	assert.True(t, item.IsGuarded)
}

func TestWhenUnguardedRuleIsApplied_ClearsItemGuarded(t *testing.T) {
	t.Parallel()
	// Arrange
	rule := content_rules.NewRuleGuarded(false)
	item := template_model.MandatoryContentItem{SID: "x", IsGuarded: true}

	// Act
	rule.Apply(&item)

	// Assert
	assert.False(t, item.IsGuarded)
}
