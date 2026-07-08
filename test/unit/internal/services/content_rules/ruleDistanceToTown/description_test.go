package ruleDistanceToTown_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/content_rules"
	"github.com/stretchr/testify/assert"
)

func TestWhenDescriptionIsQueried_ExplainsTownDistance(t *testing.T) {
	// Arrange
	rule := content_rules.NewRuleDistanceToTown(nil)

	// Act
	description := rule.Description()

	// Assert
	assert.Equal(t, "Distance to the nearest town from the content item.", description)
}
