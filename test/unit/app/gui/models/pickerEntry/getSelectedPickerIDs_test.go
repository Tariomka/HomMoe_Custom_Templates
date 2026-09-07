package pickerEntry_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/models"
	"github.com/stretchr/testify/assert"
)

func TestWhenEntriesAreSelected_TheIdsComeBackInDisplayOrder(t *testing.T) {
	t.Parallel()
	// Arrange
	entries := []models.PickerEntry{{ID: "a"}, {ID: "b"}, {ID: "c"}}

	// Act
	ids := models.GetSelectedPickerIDs(entries, map[string]bool{"c": true, "a": true})

	// Assert
	assert.Equal(t, []string{"a", "c"}, ids)
}

func TestWhenNothingIsSelected_NoIdsComeBack(t *testing.T) {
	t.Parallel()
	// Arrange
	entries := []models.PickerEntry{{ID: "a"}}

	// Act
	ids := models.GetSelectedPickerIDs(entries, map[string]bool{})

	// Assert
	assert.Empty(t, ids)
}
