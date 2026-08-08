package test_helpers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/stretchr/testify/mock"
)

// PickerEntryServiceMock is a testify mock of pickers.IPickerEntryService.
type PickerEntryServiceMock struct {
	mock.Mock
}

func (this *PickerEntryServiceMock) BuildItemPickerEntries(
	items []dtos.PickerItemDto) []dtos.PickerEntryDto {
	arguments := this.Called(items)
	entries, _ := arguments.Get(0).([]dtos.PickerEntryDto)
	return entries
}

func (this *PickerEntryServiceMock) BuildSpellPickerEntries(
	spells []dtos.PickerSpellDto) []dtos.PickerEntryDto {
	arguments := this.Called(spells)
	entries, _ := arguments.Get(0).([]dtos.PickerEntryDto)
	return entries
}

func (this *PickerEntryServiceMock) BuildValueOverridePickerEntries(
	sids []string) []dtos.PickerEntryDto {
	arguments := this.Called(sids)
	entries, _ := arguments.Get(0).([]dtos.PickerEntryDto)
	return entries
}

func (this *PickerEntryServiceMock) NormalizePickerFilter(text string) string {
	arguments := this.Called(text)
	return arguments.String(0)
}

func (this *PickerEntryServiceMock) GetVisiblePickerRows(
	entries []dtos.PickerEntryDto,
	filter string,
	grouped bool) []dtos.PickerRowDto {
	arguments := this.Called(entries, filter, grouped)
	rows, _ := arguments.Get(0).([]dtos.PickerRowDto)
	return rows
}

func (this *PickerEntryServiceMock) GetSelectedPickerIDs(
	entries []dtos.PickerEntryDto,
	selected map[string]bool) []string {
	arguments := this.Called(entries, selected)
	ids, _ := arguments.Get(0).([]string)
	return ids
}
