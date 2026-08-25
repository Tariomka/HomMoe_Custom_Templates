package editor_state_dto

// EditorStateDto is the editor state as it crosses between app/ and
// internal/handlers/. It is deliberately behaviour-free: no methods, no json
// tags and no schema version, because none of those are a consumer's concern.
//
// The nine groups are embedded anonymously so every field stays promoted and
// reads the same as it does on the model (dto.MapSize, dto.TemplateName, ...).
// Each group mirrors its entity counterpart exactly, which lets
// mappers.EditorStateMapper convert them as structs - Go ignores struct tags
// when deciding conversion identity, so a field added to a group but not
// mirrored here fails to compile instead of silently vanishing at the crossing.
type EditorStateDto struct {
	TemplateIdentityDto
	MapSettingsDto
	PlayerSettingsDto
	NeutralZoneSettingsDto
	CastleSettingsDto
	GenerationSettingsDto
	GameRuleSettingsDto
	ContentSettingsDto
	ManualEditSettingsDto
}
