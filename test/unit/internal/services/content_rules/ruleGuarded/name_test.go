package ruleGuarded_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/content_rules"
	"github.com/stretchr/testify/assert"
)

func TestWhenNameIsQueried_ReturnsGuarded(t *testing.T) {
	t.Parallel()
	// Arrange
	rule := content_rules.NewRuleGuarded(true)

	// Act
	name := rule.Name()

	// Assert
	assert.Equal(t, "Guarded", name)
}
