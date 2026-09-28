package dtos

import "github.com/Tariomka/hommoe_custom_templates/internal/models/content_rule_model"

type ContentRuleOptionDto struct {
	content_rule_model.ContentRuleOption

	Description string
	Marker      string
	EditorKind  content_rule_model.ContentRuleEditorKind
	EditorLabel string
}
