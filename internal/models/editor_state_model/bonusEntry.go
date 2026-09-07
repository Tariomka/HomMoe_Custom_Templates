package editor_state_model

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/linq"
)

type BonusEntry struct {
	editor_state.BonusEntry
}

// ToBonusEntryModels wraps persisted bonuses for use at the service layer.
func ToBonusEntryModels(bonuses []editor_state.BonusEntry) []BonusEntry {
	if len(bonuses) == 0 {
		return nil
	}

	return linq.FromSlice(bonuses).
		Select(func(bonus editor_state.BonusEntry) BonusEntry { return BonusEntry{BonusEntry: bonus} }).
		ToSlice()
}

// ToBonusEntryEntities unwraps bonuses back into their persisted form.
func ToBonusEntryEntities(bonuses []BonusEntry) []editor_state.BonusEntry {
	if len(bonuses) == 0 {
		return nil
	}

	return linq.FromSlice(bonuses).
		Select(func(bonus BonusEntry) editor_state.BonusEntry { return bonus.BonusEntry }).
		ToSlice()
}
