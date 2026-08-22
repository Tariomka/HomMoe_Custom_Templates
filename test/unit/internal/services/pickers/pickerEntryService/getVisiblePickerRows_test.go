package pickerEntryService_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/pickers"
	"github.com/stretchr/testify/assert"
)

func TestWhenThePickerIsGrouped_EachGroupIsHeadedOnceBeforeItsEntries(t *testing.T) {
	t.Parallel()
	// Arrange
	service := pickers.NewPickerEntryService()

	// Act
	rows := service.GetVisiblePickerRows(pickerRowFixtureEntries(), "", true)

	// Assert
	assert.Equal(t, []dtos.PickerRowDto{
		{IsGroupHeader: true, Group: "Weapons", GroupMatchCount: 2},
		{Entry: pickerRowFixtureEntries()[0]},
		{Entry: pickerRowFixtureEntries()[1]},
		{IsGroupHeader: true, Group: "Armor", GroupMatchCount: 1},
		{Entry: pickerRowFixtureEntries()[2]},
	}, rows)
}

func TestWhenThePickerIsFlat_NoGroupHeadersAreEmitted(t *testing.T) {
	t.Parallel()
	// Arrange
	service := pickers.NewPickerEntryService()

	// Act
	rows := service.GetVisiblePickerRows(pickerRowFixtureEntries(), "", false)

	// Assert
	assert.Len(t, rows, 3)
}

func TestWhenAFilterIsApplied_OnlyMatchingEntriesAndTheirGroupsRemain(t *testing.T) {
	t.Parallel()
	// Arrange
	service := pickers.NewPickerEntryService()

	// Act
	rows := service.GetVisiblePickerRows(pickerRowFixtureEntries(), "sword", true)

	// Assert
	assert.Equal(t, []dtos.PickerRowDto{
		{IsGroupHeader: true, Group: "Weapons", GroupMatchCount: 1},
		{Entry: pickerRowFixtureEntries()[0]},
	}, rows)
}

func TestWhenNothingMatchesTheFilter_NoRowsAreProduced(t *testing.T) {
	t.Parallel()
	// Arrange
	service := pickers.NewPickerEntryService()

	// Act
	rows := service.GetVisiblePickerRows(pickerRowFixtureEntries(), "nothing", true)

	// Assert
	assert.Empty(t, rows)
}

func pickerRowFixtureEntries() []dtos.PickerEntryDto {
	return []dtos.PickerEntryDto{
		{ID: "a", Group: "Weapons", Haystack: "sword"},
		{ID: "b", Group: "Weapons", Haystack: "axe"},
		{ID: "c", Group: "Armor", Haystack: "shield"},
	}
}
