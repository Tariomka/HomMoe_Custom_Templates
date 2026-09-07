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

// CloneZoneContentRow deep-clones a row.
func CloneZoneContentRow(row editor_state.ZoneContentRow) editor_state.ZoneContentRow {
	clone := row
	clone.Rules = CloneContentRuleRows(row.Rules)
	return clone
}

// CloneZoneContentRows deep-clones a row slice.
func CloneZoneContentRows(rows []editor_state.ZoneContentRow) []editor_state.ZoneContentRow {
	return linq.SelectSlice(rows, CloneZoneContentRow)
}
