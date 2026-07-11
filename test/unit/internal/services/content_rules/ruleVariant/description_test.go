package ruleVariant_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/content_rules"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenDescriptionIsQueried_ExplainsVariantForcing(t *testing.T) {
	t.Parallel()
	// Arrange
	rule, err := content_rules.NewRuleVariant(nil, nil)
	require.NoError(t, err)

	// Act
	description := rule.Description()

	// Assert
	assert.Equal(t, "Forces the content item to spawn a specific variant.", description)
}
