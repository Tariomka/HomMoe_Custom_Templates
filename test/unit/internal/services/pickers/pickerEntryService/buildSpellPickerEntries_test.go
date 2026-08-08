package pickerEntryService_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/pickers"
	"github.com/stretchr/testify/assert"
)

func TestWhenTheSchoolHasADisplayName_ItBecomesTheEntryGroup(t *testing.T) {
	t.Parallel()
	// Arrange
	service := pickers.NewPickerEntryService()
	spells := []dtos.PickerSpellDto{{
		Sid:               "Sid.Bless",
		Name:              "Bless",
		School:            "light",
		SchoolDisplayName: "Light Magic",
		Tier:              3,
	}}

	// Act
	entries := service.BuildSpellPickerEntries(spells)

	// Assert
	assert.Equal(t, []dtos.PickerEntryDto{{
		ID:       "Sid.Bless",
		Group:    "Light Magic",
		Label:    "Bless",
		Badge:    "[T3]",
		Haystack: "bless sid.bless light",
	}}, entries)
}

func TestWhenTheSchoolHasNoDisplayName_TheRawSchoolIsUsedAsTheGroup(t *testing.T) {
	t.Parallel()
	// Arrange
	service := pickers.NewPickerEntryService()
	spells := []dtos.PickerSpellDto{{Sid: "Sid.Bless", Name: "Bless", School: "light", Tier: 1}}

	// Act
	entries := service.BuildSpellPickerEntries(spells)

	// Assert
	assert.Equal(t, "light", entries[0].Group)
}

func TestWhenThereAreNoSpells_NoEntriesAreBuilt(t *testing.T) {
	t.Parallel()
	// Arrange
	service := pickers.NewPickerEntryService()

	// Act
	entries := service.BuildSpellPickerEntries(nil)

	// Assert
	assert.Empty(t, entries)
}
