package ruleVariant_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/content_rules"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenRuleIsApplied_SetsItemVariantId(t *testing.T) {
	// Arrange
	variantId := 3
	rule, err := content_rules.NewRuleVariant(&content_rules.UtopiaVariants, &variantId)
	require.NoError(t, err)
	item := entities.MandatoryContentItem{SID: "dragon_utopia"}

	// Act
	rule.Apply(&item)

	// Assert
	require.NotNil(t, item.Variant)
	assert.Equal(t, 3, *item.Variant)
}
