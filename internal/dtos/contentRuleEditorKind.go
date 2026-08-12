package dtos

type ContentRuleEditorKind string

const (
	ContentRuleEditorKindDistance ContentRuleEditorKind = "distance"
	ContentRuleEditorKindBoolean  ContentRuleEditorKind = "boolean"
	ContentRuleEditorKindVariant  ContentRuleEditorKind = "variant"
)
