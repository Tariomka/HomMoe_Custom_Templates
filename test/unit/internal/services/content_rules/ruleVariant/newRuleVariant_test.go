package ruleVariant_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/content_rules"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenMappingIsNil_DefaultsToUtopiaMapping(t *testing.T) {
	t.Parallel()
	// Arrange

	// Act
	rule, err := content_rules.NewRuleVariant(nil, nil)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, content_rules.NewVariantMappingCatalog().GetDefaultMapping(), rule.Mapping)
}

func TestWhenVariantIdIsNil_UsesSmallestDefinedId(t *testing.T) {
	t.Parallel()
	// Arrange
	mapping := models.NewVariantMapping(models.SidMapping{Sid: "x"}, []data.Tuple[int, string]{
		data.NewTuple(5, "Five"),
		data.NewTuple(2, "Two"),
		data.NewTuple(9, "Nine"),
	})

	// Act
	rule, err := content_rules.NewRuleVariant(&mapping, nil)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, 2, rule.VariantID)
}

func TestWhenVariantIdIsSupplied_StoresIt(t *testing.T) {
	t.Parallel()
	// Arrange
	variantID := 2
	defaultMapping := content_rules.NewVariantMappingCatalog().GetDefaultMapping()

	// Act
	rule, err := content_rules.NewRuleVariant(&defaultMapping, &variantID)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, 2, rule.VariantID)
}

func TestWhenVariantIdIsNotInMapping_ReturnsError(t *testing.T) {
	t.Parallel()
	// Arrange
	variantID := 99
	defaultMapping := content_rules.NewVariantMappingCatalog().GetDefaultMapping()

	// Act
	_, err := content_rules.NewRuleVariant(&defaultMapping, &variantID)

	// Assert
	assert.Error(t, err)
}
