package ruleDistanceToTown_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/content_rules"
	"github.com/stretchr/testify/assert"
)

func TestWhenMarkerIsQueried_ReturnsT(t *testing.T) {
	t.Parallel()
	// Arrange
	rule := content_rules.NewRuleDistanceToTown(nil)

	// Act
	marker := rule.Marker()

	// Assert
	assert.Equal(t, "T", marker)
}
