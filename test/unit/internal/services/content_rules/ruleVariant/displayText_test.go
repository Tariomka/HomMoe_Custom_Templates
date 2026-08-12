package ruleVariant_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/content_rules"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenVariantIdIsKnown_ShowsVariantDescription(t *testing.T) {
	t.Parallel()
	// Arrange
	variantID := 2
	defaultMapping := content_rules.NewVariantMappingCatalog().GetDefaultMapping()
	rule, err := content_rules.NewRuleVariant(&defaultMapping, &variantID)
	require.NoError(t, err)

	// Act
	displayText := rule.DisplayText()

	// Assert
	assert.Equal(t, "Variant: Large Guard", displayText)
}

// The error branch is only reachable by constructing the rule directly with a
// mismatched id, since the constructor validates the id against the mapping.
func TestWhenVariantIdIsUnknown_ShowsUnforeseenError(t *testing.T) {
	t.Parallel()
	// Arrange
	defaultMapping := content_rules.NewVariantMappingCatalog().GetDefaultMapping()
	rule := content_rules.RuleVariant{Mapping: defaultMapping, VariantID: 99}

	// Act
	displayText := rule.DisplayText()

	// Assert
	assert.Equal(t, "Variant: Unforeseen Error", displayText)
}
