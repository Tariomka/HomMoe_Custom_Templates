package editor_state_model

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/editor_state_helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/linq"
)

// ZoneContentRow adds behaviour to the behaviour-free
// editor_state.ZoneContentRow entity.
type ZoneContentRow struct {
	editor_state.ZoneContentRow

	// Rules shadows the embedded entity slice so the rules are model types too.
	// The embedded slice is left nil; this field is the only source of truth.
	Rules []ContentRuleRow
}

// ToZoneContentRowModels wraps persisted rows for use at the service layer.
func ToZoneContentRowModels(rows []editor_state.ZoneContentRow) []ZoneContentRow {
	if len(rows) == 0 {
		return nil
	}

	return linq.FromSlice(rows).
		Select(func(row editor_state.ZoneContentRow) ZoneContentRow {
			rules := ToContentRuleRowModels(row.Rules)
			row.Rules = nil
			return ZoneContentRow{ZoneContentRow: row, Rules: rules}
		}).ToSlice()
}

// ToZoneContentRowEntities unwraps rows back into their persisted form, folding
// the model rules back into the entity they were lifted out of.
func ToZoneContentRowEntities(rows []ZoneContentRow) []editor_state.ZoneContentRow {
	if len(rows) == 0 {
		return nil
	}

	return linq.FromSlice(rows).
		Select(func(row ZoneContentRow) editor_state.ZoneContentRow {
			entity := row.ZoneContentRow
			entity.Rules = ToContentRuleRowEntities(row.Rules)
			return entity
		}).ToSlice()
}

// CloneZoneContentRows deep-clones a row slice.
func CloneZoneContentRows(rows []ZoneContentRow) []ZoneContentRow {
	return linq.FromSlice(rows).Select(ZoneContentRow.Clone).ToSlice()
}

// Clone returns a copy sharing no backing array or pointer with the receiver.
func (this ZoneContentRow) Clone() ZoneContentRow {
	return ZoneContentRow{
		ZoneContentRow: editor_state_helpers.CloneZoneContentRow(this.ZoneContentRow),
		Rules:          CloneContentRuleRows(this.Rules),
	}
}

// Normalized returns a copy with the default values applied.
func (this ZoneContentRow) Normalized() ZoneContentRow {
	this.ZoneContentRow = editor_state_helpers.NormalizeZoneContentRow(this.ZoneContentRow)
	return this
}
