package pickers

import (
	"fmt"
	"strings"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
)

type PickerEntryService struct{}

func NewPickerEntryService() IPickerEntryService {
	return &PickerEntryService{}
}

func (this *PickerEntryService) BuildItemPickerEntries(
	items []dtos.PickerItemDto) []dtos.PickerEntryDto {
	entries := make([]dtos.PickerEntryDto, 0, len(items))
	for _, item := range items {
		entries = append(entries, dtos.PickerEntryDto{
			ID:       item.Sid,
			Group:    item.Category,
			Label:    item.Name,
			Haystack: strings.ToLower(item.Name + " " + item.Sid + " " + item.Category),
		})
	}

	return entries
}

func (this *PickerEntryService) BuildSpellPickerEntries(
	spells []dtos.PickerSpellDto) []dtos.PickerEntryDto {
	entries := make([]dtos.PickerEntryDto, 0, len(spells))
	for _, spell := range spells {
		group := spell.SchoolDisplayName
		if group == "" {
			group = spell.School
		}
		entries = append(entries, dtos.PickerEntryDto{
			ID:       spell.Sid,
			Group:    group,
			Label:    spell.Name,
			Badge:    fmt.Sprintf("[T%d]", spell.Tier),
			Haystack: strings.ToLower(spell.Name + " " + spell.Sid + " " + spell.School),
		})
	}

	return entries
}

func (this *PickerEntryService) BuildValueOverridePickerEntries(
	sids []string) []dtos.PickerEntryDto {
	entries := make([]dtos.PickerEntryDto, 0, len(sids))
	for _, sid := range sids {
		entries = append(entries, dtos.PickerEntryDto{
			ID:       sid,
			Label:    sid,
			Haystack: strings.ToLower(sid),
		})
	}

	return entries
}

func (this *PickerEntryService) NormalizePickerFilter(text string) string {
	return strings.ToLower(strings.TrimSpace(text))
}

// GetVisiblePickerRows flattens the filtered entries into group headers and leaf
// rows, preserving the entry order so callers keep control of grouping/sorting.
func (this *PickerEntryService) GetVisiblePickerRows(
	entries []dtos.PickerEntryDto,
	filter string,
	grouped bool) []dtos.PickerRowDto {
	var rows []dtos.PickerRowDto
	emitted := map[string]bool{}

	for _, entry := range entries {
		if !strings.Contains(entry.Haystack, filter) {
			continue
		}
		if grouped && !emitted[entry.Group] {
			emitted[entry.Group] = true
			rows = append(rows, dtos.PickerRowDto{
				IsGroupHeader:   true,
				Group:           entry.Group,
				GroupMatchCount: countGroupMatches(entries, entry.Group, filter),
			})
		}
		rows = append(rows, dtos.PickerRowDto{Entry: entry})
	}

	return rows
}

func (this *PickerEntryService) GetSelectedPickerIDs(
	entries []dtos.PickerEntryDto,
	selected map[string]bool) []string {
	var ids []string
	for _, entry := range entries {
		if selected[entry.ID] {
			ids = append(ids, entry.ID)
		}
	}

	return ids
}

func countGroupMatches(entries []dtos.PickerEntryDto, group string, filter string) int {
	count := 0
	for _, entry := range entries {
		if entry.Group == group && strings.Contains(entry.Haystack, filter) {
			count++
		}
	}

	return count
}
