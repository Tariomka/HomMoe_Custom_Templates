package ruleVariant_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/content_rules"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenRuleIsApplied_SetsItemVariantId(t *testing.T) {
	t.Parallel()
	// Arrange
	variantID := 3
	defaultMapping := content_rules.NewVariantMappingCatalog().GetDefaultMapping()
	rule, err := content_rules.NewRuleVariant(&defaultMapping, &variantID)
	require.NoError(t, err)
	item := template_model.MandatoryContentItem{SID: "dragon_utopia"}

	// Act
	rule.Apply(&item)

	// Assert
	require.NotNil(t, item.Variant)
	assert.Equal(t, 3, *item.Variant)
}
