package pickerEntryService_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/pickers"
	"github.com/stretchr/testify/assert"
)

func TestWhenItemEntriesAreBuilt_TheItemIsMappedOntoAnEntry(t *testing.T) {
	t.Parallel()
	// Arrange
	service := pickers.NewPickerEntryService()
	items := []dtos.PickerItemDto{{Sid: "Sid.Sword", Name: "Sword", Category: "Weapons"}}

	// Act
	entries := service.BuildItemPickerEntries(items)

	// Assert
	assert.Equal(t, []dtos.PickerEntryDto{{
		ID:       "Sid.Sword",
		Group:    "Weapons",
		Label:    "Sword",
		Haystack: "sword sid.sword weapons",
	}}, entries)
}

func TestWhenThereAreNoItems_NoEntriesAreBuilt(t *testing.T) {
	t.Parallel()
	// Arrange
	service := pickers.NewPickerEntryService()

	// Act
	entries := service.BuildItemPickerEntries(nil)

	// Assert
	assert.Empty(t, entries)
}
