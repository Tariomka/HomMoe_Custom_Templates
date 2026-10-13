package valueOverrideEntry_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/stretchr/testify/assert"
)

func TestWhenEntriesAreWrapped_OrderAndDuplicatesAreKept(t *testing.T) {
	t.Parallel()
	// Arrange
	entities := []editor_state.ValueOverrideEntry{
		{Sid: "b", GuardValue: 1},
		{Sid: "a", GuardValue: 2},
		{Sid: "b", GuardValue: 3},
	}

	// Act
	models := editor_state_model.ToValueOverrideEntryModels(entities)

	// Assert
	assert.Equal(t, entities, editor_state_model.ToValueOverrideEntryEntities(models))
}

func TestWhenThereAreNoEntries_WrappingReturnsNil(t *testing.T) {
	t.Parallel()
	// Arrange
	entities := []editor_state.ValueOverrideEntry{}

	// Act
	models := editor_state_model.ToValueOverrideEntryModels(entities)

	// Assert
	assert.Nil(t, models)
}
