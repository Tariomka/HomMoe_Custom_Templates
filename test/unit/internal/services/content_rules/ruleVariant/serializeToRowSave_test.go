package ruleVariant_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/content_rules"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenRuleIsSerialized_WritesNameAndVariantId(t *testing.T) {
	t.Parallel()
	// Arrange
	variantID := 2
	defaultMapping := content_rules.NewVariantMappingCatalog().GetDefaultMapping()
	rule, err := content_rules.NewRuleVariant(&defaultMapping, &variantID)
	require.NoError(t, err)
	expectedID := 2

	// Act
	saved := rule.SerializeToRowSave()

	// Assert
	assert.Equal(t, editor_state_model.ContentRuleRow{Name: "Variant", VariantID: &expectedID}, saved)
}
