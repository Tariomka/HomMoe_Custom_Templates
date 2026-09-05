package zoneContentRow_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenThereAreNoRows_ReturnsNil(t *testing.T) {
	t.Parallel()
	// Arrange
	var rows []editor_state.ZoneContentRow

	// Act
	models := editor_state_model.ToZoneContentRowModels(rows)

	// Assert
	assert.Nil(t, models)
}

func TestWhenRowsArePersisted_TheScalarFieldsAreCarried(t *testing.T) {
	t.Parallel()
	// Arrange
	rows := []editor_state.ZoneContentRow{{Sid: "wood_mine", Count: 3, IsGroup: true, IsMine: true}}

	// Act
	models := editor_state_model.ToZoneContentRowModels(rows)

	// Assert
	require.Len(t, models, 1)
	assert.Equal(t, rows[0], models[0].ZoneContentRow)
}

// The model shadows the entity's Rules so the rules are model types too. The
// embedded slice must be emptied, otherwise the row would carry two rule lists
// and a later unwrap could resurrect the stale one.
func TestWhenARowCarriesRules_TheyMoveOntoTheOuterFieldAndTheEmbeddedOneIsCleared(t *testing.T) {
	t.Parallel()
	// Arrange
	rows := []editor_state.ZoneContentRow{{
		Sid:   "pandora_box",
		Rules: []editor_state.ContentRuleRow{{Name: "Guarded"}},
	}}

	// Act
	models := editor_state_model.ToZoneContentRowModels(rows)

	// Assert
	require.Len(t, models, 1)
	assert.Nil(t, models[0].ZoneContentRow.Rules)
}

func TestWhenARowCarriesRules_TheOuterFieldHoldsThem(t *testing.T) {
	t.Parallel()
	// Arrange
	rows := []editor_state.ZoneContentRow{{
		Sid:   "pandora_box",
		Rules: []editor_state.ContentRuleRow{{Name: "Guarded", DistanceName: "Near"}},
	}}

	// Act
	models := editor_state_model.ToZoneContentRowModels(rows)

	// Assert
	require.Len(t, models, 1)
	assert.Equal(
		t,
		[]editor_state_model.ContentRuleRow{{ContentRuleRow: rows[0].Rules[0]}},
		models[0].Rules)
}

// The source row is passed by value, so clearing its embedded rules must not
// reach back into the caller's slice.
func TestWhenRowsAreWrapped_TheSourceRulesAreLeftIntact(t *testing.T) {
	t.Parallel()
	// Arrange
	rows := []editor_state.ZoneContentRow{{
		Sid:   "pandora_box",
		Rules: []editor_state.ContentRuleRow{{Name: "Guarded"}},
	}}

	// Act
	_ = editor_state_model.ToZoneContentRowModels(rows)

	// Assert
	assert.Len(t, rows[0].Rules, 1)
}
