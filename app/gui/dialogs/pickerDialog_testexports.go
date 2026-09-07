//go:build integration_test

package dialogs

import (
	"strings"

	"gioui.org/widget/material"
)

// EntryIDs returns every entry the picker was built with, in display order.
// ONLY FOR INTEGRATION TEST USE
func (this *MultiSelectPicker) EntryIDs() []string {
	ids := make([]string, 0, len(this.entries))
	for _, entry := range this.entries {
		ids = append(ids, entry.ID)
	}

	return ids
}

// MatchingEntryIDs returns the entries that survive the current search filter,
// in display order. ONLY FOR INTEGRATION TEST USE
func (this *MultiSelectPicker) MatchingEntryIDs() []string {
	filter := strings.ToLower(strings.TrimSpace(this.search.Text()))
	ids := make([]string, 0, len(this.entries))
	for _, entry := range this.entries {
		if strings.Contains(entry.Haystack, filter) {
			ids = append(ids, entry.ID)
		}
	}

	return ids
}

// RowCount reports how many rows - group headers included - the current filter
// produces. ONLY FOR INTEGRATION TEST USE
func (this *MultiSelectPicker) RowCount(theme *material.Theme) int {
	filter := strings.ToLower(strings.TrimSpace(this.search.Text()))
	return len(this.getRowWidgets(theme, filter))
}

// ClickEntry queues a click on the row for id and reports whether the picker
// holds such an entry. The row must be laid out in the next frame for the click
// to land, so narrow the list with SetSearch first. ONLY FOR INTEGRATION TEST USE
func (this *MultiSelectPicker) ClickEntry(id string) bool {
	for _, entry := range this.entries {
		if entry.ID == id {
			this.clickFor(id).Click()
			return true
		}
	}

	return false
}

// SetSearch ONLY FOR INTEGRATION TEST USE
func (this *MultiSelectPicker) SetSearch(text string) { this.search.SetText(text) }

// ClickAdd ONLY FOR INTEGRATION TEST USE
func (this *MultiSelectPicker) ClickAdd() { this.addBtn.Click() }

// ClickCancel ONLY FOR INTEGRATION TEST USE
func (this *MultiSelectPicker) ClickCancel() { this.cancelBtn.Click() }
