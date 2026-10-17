package banEntry_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/stretchr/testify/assert"
)

func TestWhenEntriesAreWrapped_OrderAndDuplicatesAreKept(t *testing.T) {
	t.Parallel()
	// Arrange
	entities := []editor_state.BanEntry{{Sid: "b"}, {Sid: "a"}, {Sid: "b"}}

	// Act
	models := editor_state_model.ToBanEntryModels(entities)

	// Assert
	assert.Equal(t, entities, editor_state_model.ToBanEntryEntities(models))
}

func TestWhenThereAreNoEntries_WrappingReturnsNil(t *testing.T) {
	t.Parallel()
	// Arrange
	entities := []editor_state.BanEntry{}

	// Act
	models := editor_state_model.ToBanEntryModels(entities)

	// Assert
	assert.Nil(t, models)
}
