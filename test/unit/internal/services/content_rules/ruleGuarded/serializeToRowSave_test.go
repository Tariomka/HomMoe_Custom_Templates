package ruleGuarded_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/content_rules"
	"github.com/stretchr/testify/assert"
)

func TestWhenGuardedRuleIsSerialized_WritesNameAndTrueState(t *testing.T) {
	t.Parallel()
	// Arrange
	rule := content_rules.NewRuleGuarded(true)
	isGuarded := true

	// Act
	saved := rule.SerializeToRowSave()

	// Assert
	assert.Equal(t, editor_state_model.ContentRuleRow{Name: "Guarded", IsGuarded: &isGuarded}, saved)
}

func TestWhenUnguardedRuleIsSerialized_WritesNameAndFalseState(t *testing.T) {
	t.Parallel()
	// Arrange
	rule := content_rules.NewRuleGuarded(false)
	isGuarded := false

	// Act
	saved := rule.SerializeToRowSave()

	// Assert
	assert.Equal(t, editor_state_model.ContentRuleRow{Name: "Guarded", IsGuarded: &isGuarded}, saved)
}
