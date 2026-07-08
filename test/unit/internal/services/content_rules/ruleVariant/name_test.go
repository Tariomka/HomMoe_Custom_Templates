package ruleVariant_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/content_rules"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenNameIsQueried_ReturnsVariant(t *testing.T) {
	// Arrange
	rule, err := content_rules.NewRuleVariant(nil, nil)
	require.NoError(t, err)

	// Act
	name := rule.Name()

	// Assert
	assert.Equal(t, "Variant", name)
}
