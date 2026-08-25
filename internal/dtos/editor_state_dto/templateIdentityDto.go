package editor_state_dto

// TemplateIdentityDto mirrors the entity group of the same name. Field names,
// types and order must match it exactly: the mapper converts the two structs
// directly, so a divergence is a compile error rather than a silent data loss.
type TemplateIdentityDto struct {
	TemplateName string
	GameMode     string
}
