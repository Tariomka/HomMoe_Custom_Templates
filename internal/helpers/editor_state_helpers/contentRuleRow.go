package editor_state_helpers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/linq"
)

// CloneContentRuleRow returns a copy that shares no pointer with the original.
func CloneContentRuleRow(rule editor_state.ContentRuleRow) editor_state.ContentRuleRow {
	clone := rule
	clone.IsGuarded = helpers.ClonePointer(rule.IsGuarded)
	clone.IsSoloEncounter = helpers.ClonePointer(rule.IsSoloEncounter)
	clone.VariantID = helpers.ClonePointer(rule.VariantID)
	return clone
}

// CloneContentRuleRows deep-clones a rule slice. A nil slice stays nil, because
// the editor-state change detection distinguishes nil from empty.
func CloneContentRuleRows(rules []editor_state.ContentRuleRow) []editor_state.ContentRuleRow {
	return linq.FromSlice(rules).
		Select(func(rule editor_state.ContentRuleRow) editor_state.ContentRuleRow { return CloneContentRuleRow(rule) }).
		ToSlice()
}
