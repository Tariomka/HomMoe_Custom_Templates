package editor_state_model

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/linq"
)

type BanEntry struct {
	editor_state.BanEntry
}

func ToBanEntryModels(entries []editor_state.BanEntry) []BanEntry {
	if len(entries) == 0 {
		return nil
	}

	return linq.FromSlice(entries).
		Select(func(entry editor_state.BanEntry) BanEntry { return BanEntry{BanEntry: entry} }).
		ToSlice()
}

func ToBanEntryEntities(entries []BanEntry) []editor_state.BanEntry {
	if len(entries) == 0 {
		return nil
	}

	return linq.FromSlice(entries).
		Select(func(entry BanEntry) editor_state.BanEntry { return entry.BanEntry }).
		ToSlice()
}
