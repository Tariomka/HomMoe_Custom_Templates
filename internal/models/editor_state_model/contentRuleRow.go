package editor_state_model

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/editor_state_helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/linq"
)

// ContentRuleRow adds behaviour to the behaviour-free
// editor_state.ContentRuleRow entity.
type ContentRuleRow struct {
	editor_state.ContentRuleRow
}

// ToContentRuleRowModels wraps persisted rules for use at the service layer.
func ToContentRuleRowModels(rules []editor_state.ContentRuleRow) []ContentRuleRow {
	if len(rules) == 0 {
		return nil
	}

	return linq.FromSlice(rules).
		Select(func(rule editor_state.ContentRuleRow) ContentRuleRow { return ContentRuleRow{ContentRuleRow: rule} }).
		ToSlice()
}

// ToContentRuleRowEntities unwraps rules back into their persisted form.
func ToContentRuleRowEntities(rules []ContentRuleRow) []editor_state.ContentRuleRow {
	if len(rules) == 0 {
		return nil
	}

	return linq.FromSlice(rules).
		Select(func(rule ContentRuleRow) editor_state.ContentRuleRow { return rule.ContentRuleRow }).
		ToSlice()
}

// CloneContentRuleRows deep-clones a rule slice.
func CloneContentRuleRows(rules []ContentRuleRow) []ContentRuleRow {
	return linq.FromSlice(rules).Select(ContentRuleRow.Clone).ToSlice()
}

// Clone returns a copy sharing no pointer with the receiver.
func (this ContentRuleRow) Clone() ContentRuleRow {
	return ContentRuleRow{ContentRuleRow: editor_state_helpers.CloneContentRuleRow(this.ContentRuleRow)}
}
