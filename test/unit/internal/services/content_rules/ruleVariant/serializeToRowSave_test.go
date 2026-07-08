package ruleVariant_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/content_rules"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenRuleIsSerialized_WritesNameAndVariantId(t *testing.T) {
	// Arrange
	variantId := 2
	rule, err := content_rules.NewRuleVariant(&content_rules.UtopiaVariants, &variantId)
	require.NoError(t, err)
	expectedId := 2

	// Act
	saved := rule.SerializeToRowSave()

	// Assert
	assert.Equal(t, models.ContentRuleRowSave{Name: "Variant", VariantId: &expectedId}, saved)
}
