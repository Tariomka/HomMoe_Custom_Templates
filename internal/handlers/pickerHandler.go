package handlers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers/handler_interfaces"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/pickers"
)

type pickerHandler struct {
	pickerEntryService pickers.IPickerEntryService
}

func NewPickerHandler(pickerEntryService pickers.IPickerEntryService) handler_interfaces.IPickerHandler {
	return &pickerHandler{pickerEntryService: pickerEntryService}
}

func (this *pickerHandler) BuildItemPickerEntries(items []dtos.PickerItemDto) []dtos.PickerEntryDto {
	return this.pickerEntryService.BuildItemPickerEntries(items)
}

func (this *pickerHandler) BuildSpellPickerEntries(spells []dtos.PickerSpellDto) []dtos.PickerEntryDto {
	return this.pickerEntryService.BuildSpellPickerEntries(spells)
}

func (this *pickerHandler) BuildValueOverridePickerEntries(sids []string) []dtos.PickerEntryDto {
	return this.pickerEntryService.BuildValueOverridePickerEntries(sids)
}

func (this *pickerHandler) NormalizePickerFilter(text string) string {
	return this.pickerEntryService.NormalizePickerFilter(text)
}

func (this *pickerHandler) GetVisiblePickerRows(
	entries []dtos.PickerEntryDto,
	filter string,
	grouped bool) []dtos.PickerRowDto {
	return this.pickerEntryService.GetVisiblePickerRows(entries, filter, grouped)
}

func (this *pickerHandler) GetSelectedPickerIDs(
	entries []dtos.PickerEntryDto,
	selected map[string]bool) []string {
	return this.pickerEntryService.GetSelectedPickerIDs(entries, selected)
}
