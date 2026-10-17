package valueOverrideEntry_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/stretchr/testify/assert"
)

func TestWhenEntriesAreUnwrapped_SidAndGuardValueAreKept(t *testing.T) {
	t.Parallel()
	// Arrange
	models := []editor_state_model.ValueOverrideEntry{
		{Sid: "a", GuardValue: 700},
	}

	// Act
	entities := editor_state_model.ToValueOverrideEntryEntities(models)

	// Assert
	assert.Equal(t, []editor_state.ValueOverrideEntry{{Sid: "a", GuardValue: 700}}, entities)
}

func TestWhenThereAreNoEntries_UnwrappingReturnsNil(t *testing.T) {
	t.Parallel()
	// Arrange
	models := []editor_state_model.ValueOverrideEntry{}

	// Act
	entities := editor_state_model.ToValueOverrideEntryEntities(models)

	// Assert
	assert.Nil(t, entities)
}
