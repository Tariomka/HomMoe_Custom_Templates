package dtos

import "github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"

type ContentRuleDescriptionDto struct {
	Key          ContentRuleKey
	DisplayText  string
	Marker       string
	VariantLabel string
	Valid        bool
	SavedRule    editor_state_model.ContentRuleRow
}
