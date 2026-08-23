package dtos

import "github.com/Tariomka/hommoe_custom_templates/internal/models"

// ContentRuleCompositionResultDto is the outcome of composing a content rule
// from the manage-rules editor state. Valid is false when the editor selection
// cannot produce a rule.
type ContentRuleCompositionResultDto struct {
	Rule  models.ContentRuleRow
	Valid bool
}
