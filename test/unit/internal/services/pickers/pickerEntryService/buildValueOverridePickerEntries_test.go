package pickerEntryService_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/pickers"
	"github.com/stretchr/testify/assert"
)

func TestWhenValueOverrideEntriesAreBuilt_TheSidIsBothIdAndLabel(t *testing.T) {
	t.Parallel()
	// Arrange
	service := pickers.NewPickerEntryService()

	// Act
	entries := service.BuildValueOverridePickerEntries([]string{"Sid.Gold"})

	// Assert
	assert.Equal(t, []dtos.PickerEntryDto{{
		ID:       "Sid.Gold",
		Label:    "Sid.Gold",
		Haystack: "sid.gold",
	}}, entries)
}

func TestWhenThereAreNoValueOverrideSids_NoEntriesAreBuilt(t *testing.T) {
	t.Parallel()
	// Arrange
	service := pickers.NewPickerEntryService()

	// Act
	entries := service.BuildValueOverridePickerEntries(nil)

	// Assert
	assert.Empty(t, entries)
}
