package zoneContentRow_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenARowIsCloned_TheRulesAreEqual(t *testing.T) {
	t.Parallel()
	// Arrange
	row := editor_state_model.ZoneContentRow{
		Sid:   "pandora_box",
		Count: 2,
		Rules: []editor_state_model.ContentRuleRow{{Name: "Guarded", IsGuarded: new(true)}},
	}

	// Act
	clone := row.Clone()

	// Assert
	assert.Equal(t, row, clone)
}

func TestWhenAClonedRowRuleIsMutated_TheSourceIsUnaffected(t *testing.T) {
	t.Parallel()
	// Arrange
	row := editor_state_model.ZoneContentRow{
		Sid:   "pandora_box",
		Rules: []editor_state_model.ContentRuleRow{{Name: "Guarded", IsGuarded: new(true)}},
	}
	clone := row.Clone()
	require.NotEmpty(t, clone.Rules)

	// Act
	clone.Rules[0].Name = "Mutated"

	// Assert
	assert.Equal(t, "Guarded", row.Rules[0].Name)
}

func TestWhenAClonedRowRulePointerIsMutated_TheSourceIsUnaffected(t *testing.T) {
	t.Parallel()
	// Arrange
	row := editor_state_model.ZoneContentRow{
		Sid:   "pandora_box",
		Rules: []editor_state_model.ContentRuleRow{{Name: "Guarded", IsGuarded: new(true)}},
	}
	clone := row.Clone()
	require.NotEmpty(t, clone.Rules)

	// Act
	*clone.Rules[0].IsGuarded = false

	// Assert
	assert.True(t, *row.Rules[0].IsGuarded)
}
