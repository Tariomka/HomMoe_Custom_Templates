package pickerEntry_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/models"
	"github.com/stretchr/testify/assert"
)

func TestWhenTheSchoolHasADisplayName_ItBecomesTheEntryGroup(t *testing.T) {
	t.Parallel()
	// Arrange
	spells := []models.PickerSpell{{
		Sid:               "Sid.Bless",
		Name:              "Bless",
		School:            "light",
		SchoolDisplayName: "Light Magic",
		Tier:              3,
	}}

	// Act
	entries := models.BuildSpellPickerEntries(spells)

	// Assert
	assert.Equal(t, []models.PickerEntry{{
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
	spells := []models.PickerSpell{{Sid: "Sid.Bless", Name: "Bless", School: "light", Tier: 1}}

	// Act
	entries := models.BuildSpellPickerEntries(spells)

	// Assert
	assert.Equal(t, "light", entries[0].Group)
}

func TestWhenThereAreNoSpells_NoEntriesAreBuilt(t *testing.T) {
	t.Parallel()
	// Arrange
	// Act
	entries := models.BuildSpellPickerEntries(nil)

	// Assert
	assert.Empty(t, entries)
}
