package contentRuleRowSave_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenRuleIsCloned_ScalarFieldsAreCopied(t *testing.T) {
	t.Parallel()
	// Arrange
	rule := models.ContentRuleRowSave{Name: gofakeit.LetterN(8), DistanceName: gofakeit.LetterN(6)}

	// Act
	clone := rule.Clone()

	// Assert
	assert.Equal(t, rule, clone)
}

func TestWhenPointersAreNil_ClonePointersStayNil(t *testing.T) {
	t.Parallel()
	// Arrange
	rule := models.ContentRuleRowSave{Name: gofakeit.LetterN(8)}

	// Act
	clone := rule.Clone()

	// Assert
	assert.Nil(t, clone.IsGuarded)
}

func TestWhenGuardedPointerIsMutatedOnTheClone_SourceIsUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	rule := models.ContentRuleRowSave{Name: gofakeit.LetterN(8), IsGuarded: new(true)}
	clone := rule.Clone()

	// Act
	*clone.IsGuarded = false

	// Assert
	require.NotNil(t, rule.IsGuarded)
	assert.True(t, *rule.IsGuarded)
}

func TestWhenSoloEncounterPointerIsMutatedOnTheClone_SourceIsUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	rule := models.ContentRuleRowSave{Name: gofakeit.LetterN(8), IsSoloEncounter: new(true)}
	clone := rule.Clone()

	// Act
	*clone.IsSoloEncounter = false

	// Assert
	require.NotNil(t, rule.IsSoloEncounter)
	assert.True(t, *rule.IsSoloEncounter)
}

func TestWhenVariantIdPointerIsMutatedOnTheClone_SourceIsUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	variantID := gofakeit.IntRange(1, 100)
	rule := models.ContentRuleRowSave{Name: gofakeit.LetterN(8), VariantID: &variantID}
	clone := rule.Clone()

	// Act
	*clone.VariantID = variantID + 1

	// Assert
	require.NotNil(t, rule.VariantID)
	assert.Equal(t, variantID, *rule.VariantID)
}
