package zoneContentRow_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenThereAreNoRowModels_ReturnsNil(t *testing.T) {
	t.Parallel()
	// Arrange
	var rows []editor_state_model.ZoneContentRow

	// Act
	entities := editor_state_model.ToZoneContentRowEntities(rows)

	// Assert
	assert.Nil(t, entities)
}

// The rules live on the outer field while the entity expects them on the
// embedded one, so unwrapping has to fold them back in.
func TestWhenARowModelCarriesRules_TheyAreFoldedBackIntoTheEntity(t *testing.T) {
	t.Parallel()
	// Arrange
	rule := editor_state.ContentRuleRow{Name: "Guarded", DistanceName: "Near"}
	rows := []editor_state_model.ZoneContentRow{{
		Sid:   "pandora_box",
		Rules: []editor_state_model.ContentRuleRow{{ContentRuleRow: rule}},
	}}

	// Act
	entities := editor_state_model.ToZoneContentRowEntities(rows)

	// Assert
	require.Len(t, entities, 1)
	assert.Equal(
		t,
		editor_state.ZoneContentRow{Sid: "pandora_box", Rules: []editor_state.ContentRuleRow{rule}},
		entities[0],
	)
}

func TestWhenRowsAreUnwrappedAndWrappedAgain_TheyRoundTrip(t *testing.T) {
	t.Parallel()
	// Arrange
	rows := []editor_state_model.ZoneContentRow{{
		Sid:     "wood_mine",
		Count:   2,
		IsGroup: true,
		IsMine:  true,
		Rules:   []editor_state_model.ContentRuleRow{{Name: "Guarded"}},
	}}

	// Act
	roundTripped := editor_state_model.ToZoneContentRowModels(editor_state_model.ToZoneContentRowEntities(rows))

	// Assert
	assert.Equal(t, rows, roundTripped)
}
