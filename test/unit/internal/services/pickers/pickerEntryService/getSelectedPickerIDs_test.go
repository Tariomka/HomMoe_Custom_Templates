package pickerEntryService_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/pickers"
	"github.com/stretchr/testify/assert"
)

func TestWhenEntriesAreSelected_TheIdsComeBackInDisplayOrder(t *testing.T) {
	t.Parallel()
	// Arrange
	service := pickers.NewPickerEntryService()
	entries := []dtos.PickerEntryDto{{ID: "a"}, {ID: "b"}, {ID: "c"}}

	// Act
	ids := service.GetSelectedPickerIDs(entries, map[string]bool{"c": true, "a": true})

	// Assert
	assert.Equal(t, []string{"a", "c"}, ids)
}

func TestWhenNothingIsSelected_NoIdsComeBack(t *testing.T) {
	t.Parallel()
	// Arrange
	service := pickers.NewPickerEntryService()
	entries := []dtos.PickerEntryDto{{ID: "a"}}

	// Act
	ids := service.GetSelectedPickerIDs(entries, map[string]bool{})

	// Assert
	assert.Empty(t, ids)
}
