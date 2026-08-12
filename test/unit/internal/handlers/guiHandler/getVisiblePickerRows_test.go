package guiHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/stretchr/testify/assert"
)

func TestWhenVisibleRowsAreRequested_TheGroupHeaderLeadsItsEntries(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	entries := []dtos.PickerEntryDto{{ID: "a", Group: "Weapons", Haystack: "sword"}}

	// Act
	rows := handler.GetVisiblePickerRows(entries, "", true)

	// Assert
	assert.Equal(t, []dtos.PickerRowDto{
		{IsGroupHeader: true, Group: "Weapons", GroupMatchCount: 1},
		{Entry: entries[0]},
	}, rows)
}
