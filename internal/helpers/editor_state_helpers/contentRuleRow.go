package editor_state_helpers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/linq"
)

// CloneContentRuleRow deep-clones a rule row.
func CloneContentRuleRow(rule editor_state.ContentRuleRow) editor_state.ContentRuleRow {
	clone := rule
	clone.IsGuarded = helpers.ClonePointer(rule.IsGuarded)
	clone.IsSoloEncounter = helpers.ClonePointer(rule.IsSoloEncounter)
	clone.VariantID = helpers.ClonePointer(rule.VariantID)
	return clone
}

// CloneContentRuleRows deep-clones a rule slice.
func CloneContentRuleRows(rules []editor_state.ContentRuleRow) []editor_state.ContentRuleRow {
	return linq.SelectSlice(rules, CloneContentRuleRow)
}
