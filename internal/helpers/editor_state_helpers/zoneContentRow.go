package editor_state_helpers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/linq"
)

// NormalizeZoneContentRow returns a copy with the default values applied.
func NormalizeZoneContentRow(row editor_state.ZoneContentRow) editor_state.ZoneContentRow {
	out := row
	if out.Count < 1 {
		out.Count = 1
	}
	return out
}

// CloneZoneContentRow returns a copy that shares no backing array or pointer
// with the original. A nil Rules slice stays nil, because the editor-state
// change detection distinguishes a nil slice from an empty one.
func CloneZoneContentRow(row editor_state.ZoneContentRow) editor_state.ZoneContentRow {
	clone := row
	clone.Rules = CloneContentRuleRows(row.Rules)
	return clone
}

// CloneZoneContentRows deep-clones a row slice, preserving nil.
func CloneZoneContentRows(rows []editor_state.ZoneContentRow) []editor_state.ZoneContentRow {
	return linq.FromSlice(rows).
		Select(func(row editor_state.ZoneContentRow) editor_state.ZoneContentRow { return CloneZoneContentRow(row) }).
		ToSlice()
}
