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
	assert.Equal(t, content_rules.UtopiaVariants, rule.Mapping)
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

	// Act
	rule, err := content_rules.NewRuleVariant(&content_rules.UtopiaVariants, &variantID)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, 2, rule.VariantID)
}

func TestWhenVariantIdIsNotInMapping_ReturnsError(t *testing.T) {
	t.Parallel()
	// Arrange
	variantID := 99

	// Act
	_, err := content_rules.NewRuleVariant(&content_rules.UtopiaVariants, &variantID)

	// Assert
	assert.Error(t, err)
}
