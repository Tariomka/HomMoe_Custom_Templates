package editor_state_model

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/linq"
)

type ValueOverrideEntry struct {
	editor_state.ValueOverrideEntry
}

func ToValueOverrideEntryModels(entries []editor_state.ValueOverrideEntry) []ValueOverrideEntry {
	if len(entries) == 0 {
		return nil
	}

	return linq.FromSlice(entries).
		Select(func(entry editor_state.ValueOverrideEntry) ValueOverrideEntry {
			return ValueOverrideEntry{ValueOverrideEntry: entry}
		}).
		ToSlice()
}

func ToValueOverrideEntryEntities(entries []ValueOverrideEntry) []editor_state.ValueOverrideEntry {
	if len(entries) == 0 {
		return nil
	}

	return linq.FromSlice(entries).
		Select(func(entry ValueOverrideEntry) editor_state.ValueOverrideEntry { return entry.ValueOverrideEntry }).
		ToSlice()
}
