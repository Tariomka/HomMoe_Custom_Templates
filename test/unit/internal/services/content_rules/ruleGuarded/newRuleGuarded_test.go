package ruleGuarded_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/content_rules"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenStateIsSupplied_StoresIt(t *testing.T) {
	// Arrange
	isGuarded := gofakeit.Bool()

	// Act
	rule := content_rules.NewRuleGuarded(isGuarded)

	// Assert
	assert.Equal(t, isGuarded, rule.IsGuarded)
}
