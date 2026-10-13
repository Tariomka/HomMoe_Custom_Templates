package banEntry_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/stretchr/testify/assert"
)

func TestWhenEntriesAreUnwrapped_TheirSidsAreKept(t *testing.T) {
	t.Parallel()
	// Arrange
	models := []editor_state_model.BanEntry{{Sid: "a"}}

	// Act
	entities := editor_state_model.ToBanEntryEntities(models)

	// Assert
	assert.Equal(t, []editor_state.BanEntry{{Sid: "a"}}, entities)
}

func TestWhenThereAreNoEntries_UnwrappingReturnsNil(t *testing.T) {
	t.Parallel()
	// Arrange
	models := []editor_state_model.BanEntry{}

	// Act
	entities := editor_state_model.ToBanEntryEntities(models)

	// Assert
	assert.Nil(t, entities)
}
