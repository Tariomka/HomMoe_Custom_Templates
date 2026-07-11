package ruleVariant_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/content_rules"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenMarkerIsQueried_ReturnsEmptyBadge(t *testing.T) {
	t.Parallel()
	// Arrange
	rule, err := content_rules.NewRuleVariant(nil, nil)
	require.NoError(t, err)

	// Act
	marker := rule.Marker()

	// Assert
	assert.Empty(t, marker)
}
