package handler_interfaces

import "github.com/Tariomka/hommoe_custom_templates/internal/dtos"

// IPickerHandler exposes the picker entry/row/selection logic to the GUI layer.
type IPickerHandler interface {
	BuildItemPickerEntries(items []dtos.PickerItemDto) []dtos.PickerEntryDto
	BuildSpellPickerEntries(spells []dtos.PickerSpellDto) []dtos.PickerEntryDto
	BuildValueOverridePickerEntries(sids []string) []dtos.PickerEntryDto
	NormalizePickerFilter(text string) string
	GetVisiblePickerRows(entries []dtos.PickerEntryDto, filter string, grouped bool) []dtos.PickerRowDto
	GetSelectedPickerIDs(entries []dtos.PickerEntryDto, selected map[string]bool) []string
}
