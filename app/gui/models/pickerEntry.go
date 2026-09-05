package models

import (
	"fmt"
	"strings"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/linq"
)

// PickerEntry is one selectable row of a multi-select picker, already mapped
// out of its source catalogue and ready to be searched and rendered.
type PickerEntry struct {
	ID       string
	Group    string // category / school; ignored when the picker is flat
	Label    string // primary display text
	Badge    string // optional leading badge (e.g. "[T3]")
	Trailing string // optional dim trailing text (e.g. the raw SID)
	Haystack string // lowercased search text
}

// BuildItemPickerEntries maps the bannable-item catalogue into picker entries.
func BuildItemPickerEntries(items []PickerItem) []PickerEntry {
	return linq.SelectSlice(items, func(item PickerItem) PickerEntry {
		return PickerEntry{
			ID:       item.Sid,
			Group:    item.Category,
			Label:    item.Name,
			Haystack: strings.ToLower(item.Name + " " + item.Sid + " " + item.Category),
		}
	})
}

// BuildSpellPickerEntries maps the spell catalogue into picker entries, falling
// back to the raw school when the catalogue has no display name for it.
func BuildSpellPickerEntries(spells []PickerSpell) []PickerEntry {
	return linq.SelectSlice(spells, func(spell PickerSpell) PickerEntry {
		group := spell.SchoolDisplayName
		if group == "" {
			group = spell.School
		}

		return PickerEntry{
			ID:       spell.Sid,
			Group:    group,
			Label:    spell.Name,
			Badge:    fmt.Sprintf("[T%d]", spell.Tier),
			Haystack: strings.ToLower(spell.Name + " " + spell.Sid + " " + spell.School),
		}
	})
}

// BuildValueOverridePickerEntries maps bare guard-value SIDs into picker
// entries; this picker is flat, so the entries carry no group.
func BuildValueOverridePickerEntries(sids []string) []PickerEntry {
	return linq.SelectSlice(sids, func(sid string) PickerEntry {
		return PickerEntry{
			ID:       sid,
			Label:    sid,
			Haystack: strings.ToLower(sid),
		}
	})
}

// NormalizePickerFilter turns raw search-box text into the form Haystack is
// matched against.
func NormalizePickerFilter(text string) string {
	return strings.ToLower(strings.TrimSpace(text))
}

// GetSelectedPickerIDs returns the ticked entry IDs in entry order, so the
// caller's selection map does not leak its iteration order into the result.
func GetSelectedPickerIDs(entries []PickerEntry, selected map[string]bool) []string {
	var ids []string
	for _, entry := range entries {
		if selected[entry.ID] {
			ids = append(ids, entry.ID)
		}
	}

	return ids
}
