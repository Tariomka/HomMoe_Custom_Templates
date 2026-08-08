package pickers

import "github.com/Tariomka/hommoe_custom_templates/internal/dtos"

// IPickerEntryService maps picker source catalogues into searchable entries and
// derives the filtered row model and the selection the picker dialogs render.
type IPickerEntryService interface {
	BuildItemPickerEntries(items []dtos.PickerItemDto) []dtos.PickerEntryDto
	BuildSpellPickerEntries(spells []dtos.PickerSpellDto) []dtos.PickerEntryDto
	BuildValueOverridePickerEntries(sids []string) []dtos.PickerEntryDto
	NormalizePickerFilter(text string) string
	GetVisiblePickerRows(entries []dtos.PickerEntryDto, filter string, grouped bool) []dtos.PickerRowDto
	GetSelectedPickerIDs(entries []dtos.PickerEntryDto, selected map[string]bool) []string
}
