package contentRuleRow_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/editor_state_helpers"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenRuleIsCloned_ScalarFieldsAreCopied(t *testing.T) {
	t.Parallel()
	// Arrange
	rule := editor_state.ContentRuleRow{Name: gofakeit.LetterN(8), DistanceName: gofakeit.LetterN(6)}

	// Act
	clone := editor_state_helpers.CloneContentRuleRow(rule)

	// Assert
	assert.Equal(t, rule, clone)
}

func TestWhenPointersAreNil_ClonePointersStayNil(t *testing.T) {
	t.Parallel()
	// Arrange
	rule := editor_state.ContentRuleRow{Name: gofakeit.LetterN(8)}

	// Act
	clone := editor_state_helpers.CloneContentRuleRow(rule)

	// Assert
	assert.Nil(t, clone.IsGuarded)
}

func TestWhenGuardedPointerIsMutatedOnTheClone_SourceIsUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	rule := editor_state.ContentRuleRow{Name: gofakeit.LetterN(8), IsGuarded: new(true)}
	clone := editor_state_helpers.CloneContentRuleRow(rule)

	// Act
	*clone.IsGuarded = false

	// Assert
	require.NotNil(t, rule.IsGuarded)
	assert.True(t, *rule.IsGuarded)
}

func TestWhenSoloEncounterPointerIsMutatedOnTheClone_SourceIsUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	rule := editor_state.ContentRuleRow{Name: gofakeit.LetterN(8), IsSoloEncounter: new(true)}
	clone := editor_state_helpers.CloneContentRuleRow(rule)

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
	rule := editor_state.ContentRuleRow{Name: gofakeit.LetterN(8), VariantID: &variantID}
	clone := editor_state_helpers.CloneContentRuleRow(rule)

	// Act
	*clone.VariantID = variantID + 1

	// Assert
	require.NotNil(t, rule.VariantID)
	assert.Equal(t, variantID, *rule.VariantID)
}
