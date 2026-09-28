package content_rule_model

import "github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"

// ContentRuleDescription is a saved content rule resolved against its content;
// Valid is false when the saved data does not match a known rule.
type ContentRuleDescription struct {
	Key          ContentRuleKey
	DisplayText  string
	Marker       string
	VariantLabel string
	Valid        bool
	SavedRule    editor_state_model.ContentRuleRow
}
