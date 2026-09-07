package pickerEntry_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/models"
	"github.com/stretchr/testify/assert"
)

func TestWhenItemEntriesAreBuilt_TheItemIsMappedOntoAnEntry(t *testing.T) {
	t.Parallel()
	// Arrange
	items := []models.PickerItem{{Sid: "Sid.Sword", Name: "Sword", Category: "Weapons"}}

	// Act
	entries := models.BuildItemPickerEntries(items)

	// Assert
	assert.Equal(t, []models.PickerEntry{{
		ID:       "Sid.Sword",
		Group:    "Weapons",
		Label:    "Sword",
		Haystack: "sword sid.sword weapons",
	}}, entries)
}

func TestWhenThereAreNoItems_NoEntriesAreBuilt(t *testing.T) {
	t.Parallel()
	// Arrange
	// Act
	entries := models.BuildItemPickerEntries(nil)

	// Assert
	assert.Empty(t, entries)
}
