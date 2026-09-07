package models

import "strings"

// PickerRow is one line of a picker's flattened, filtered row model: either a
// group header (with the number of matching entries below it) or a leaf entry.
type PickerRow struct {
	IsGroupHeader   bool
	Group           string
	GroupMatchCount int
	Entry           PickerEntry
}

// GetVisiblePickerRows flattens the filtered entries into group headers and leaf
// rows, preserving the entry order so callers keep control of grouping/sorting.
func GetVisiblePickerRows(entries []PickerEntry, filter string, grouped bool) []PickerRow {
	var rows []PickerRow
	emitted := map[string]bool{}

	for _, entry := range entries {
		if !strings.Contains(entry.Haystack, filter) {
			continue
		}
		if grouped && !emitted[entry.Group] {
			emitted[entry.Group] = true
			rows = append(rows, PickerRow{
				IsGroupHeader:   true,
				Group:           entry.Group,
				GroupMatchCount: countGroupMatches(entries, entry.Group, filter),
			})
		}
		rows = append(rows, PickerRow{Entry: entry})
	}

	return rows
}

func countGroupMatches(entries []PickerEntry, group string, filter string) int {
	count := 0
	for _, entry := range entries {
		if entry.Group == group && strings.Contains(entry.Haystack, filter) {
			count++
		}
	}

	return count
}
