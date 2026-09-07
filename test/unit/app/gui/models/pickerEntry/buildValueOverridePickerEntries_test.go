package pickerEntry_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/models"
	"github.com/stretchr/testify/assert"
)

func TestWhenValueOverrideEntriesAreBuilt_TheSidIsBothIdAndLabel(t *testing.T) {
	t.Parallel()
	// Arrange
	// Act
	entries := models.BuildValueOverridePickerEntries([]string{"Sid.Gold"})

	// Assert
	assert.Equal(t, []models.PickerEntry{{
		ID:       "Sid.Gold",
		Label:    "Sid.Gold",
		Haystack: "sid.gold",
	}}, entries)
}

func TestWhenThereAreNoValueOverrideSids_NoEntriesAreBuilt(t *testing.T) {
	t.Parallel()
	// Arrange
	// Act
	entries := models.BuildValueOverridePickerEntries(nil)

	// Assert
	assert.Empty(t, entries)
}
