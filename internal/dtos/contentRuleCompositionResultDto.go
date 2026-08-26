package dtos

import "github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"

type ContentRuleCompositionResultDto struct {
	Rule  editor_state_model.ContentRuleRow
	Valid bool
}
