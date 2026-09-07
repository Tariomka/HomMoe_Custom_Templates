package zoneContentRow_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenACloneOfARowSliceIsMutated_TheSourceIsUnaffected(t *testing.T) {
	t.Parallel()
	// Arrange
	rows := []editor_state_model.ZoneContentRow{{
		Sid:   "pandora_box",
		Rules: []editor_state_model.ContentRuleRow{{Name: "Guarded"}},
	}}
	clone := editor_state_model.CloneZoneContentRows(rows)
	require.Len(t, clone, 1)

	// Act
	clone[0].Rules[0].Name = "Mutated"

	// Assert
	assert.Equal(t, "Guarded", rows[0].Rules[0].Name)
}

func TestWhenARowSliceIsCloned_TheContentsAreEqual(t *testing.T) {
	t.Parallel()
	// Arrange
	rows := []editor_state_model.ZoneContentRow{{Sid: "wood_mine", Count: 2}}

	// Act
	clone := editor_state_model.CloneZoneContentRows(rows)

	// Assert
	assert.Equal(t, rows, clone)
}
