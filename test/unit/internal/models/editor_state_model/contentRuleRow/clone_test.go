package contentRuleRow_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenARuleIsCloned_TheValuesAreEqual(t *testing.T) {
	t.Parallel()
	// Arrange
	rule := editor_state_model.ContentRuleRow{
		Name:            "Guarded",
		DistanceName:    "Near",
		IsGuarded:       new(true),
		IsSoloEncounter: new(false),
		VariantID:       new(7),
	}

	// Act
	clone := rule.Clone()

	// Assert
	assert.Equal(t, rule, clone)
}

func TestWhenAClonedRulePointerIsMutated_TheSourceIsUnaffected(t *testing.T) {
	t.Parallel()
	// Arrange
	rule := editor_state_model.ContentRuleRow{Name: "Guarded", IsGuarded: new(true)}
	clone := rule.Clone()
	require.NotNil(t, clone.IsGuarded)

	// Act
	*clone.IsGuarded = false

	// Assert
	assert.True(t, *rule.IsGuarded)
}

func TestWhenACloneOfARuleSliceIsMutated_TheSourceIsUnaffected(t *testing.T) {
	t.Parallel()
	// Arrange
	rules := []editor_state_model.ContentRuleRow{{Name: "Guarded", VariantID: new(3)}}
	clone := editor_state_model.CloneContentRuleRows(rules)
	require.Len(t, clone, 1)

	// Act
	*clone[0].VariantID = 9

	// Assert
	assert.Equal(t, 3, *rules[0].VariantID)
}
