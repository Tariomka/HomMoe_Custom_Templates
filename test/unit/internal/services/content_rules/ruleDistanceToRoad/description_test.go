package ruleDistanceToRoad_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/content_rules"
	"github.com/stretchr/testify/assert"
)

func TestWhenDescriptionIsQueried_ExplainsRoadDistance(t *testing.T) {
	t.Parallel()
	// Arrange
	rule := content_rules.NewRuleDistanceToRoad(nil)

	// Act
	description := rule.Description()

	// Assert
	assert.Equal(t, "Distance to the nearest road from the content item.", description)
}
