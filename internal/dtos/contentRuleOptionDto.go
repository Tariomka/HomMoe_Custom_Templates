package dtos

type ContentRuleOptionDto struct {
	Key         ContentRuleKey
	Name        string
	Description string
	Marker      string
	EditorKind  ContentRuleEditorKind
	EditorLabel string
}
