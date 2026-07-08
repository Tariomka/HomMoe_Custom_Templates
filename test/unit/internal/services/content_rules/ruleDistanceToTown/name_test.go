package ruleDistanceToTown_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/content_rules"
	"github.com/stretchr/testify/assert"
)

func TestWhenNameIsQueried_ReturnsDistanceToTown(t *testing.T) {
	// Arrange
	rule := content_rules.NewRuleDistanceToTown(nil)

	// Act
	name := rule.Name()

	// Assert
	assert.Equal(t, "Distance to town", name)
}
