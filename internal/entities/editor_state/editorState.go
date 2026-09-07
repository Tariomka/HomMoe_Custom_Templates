package editor_state

// CurrentEditorStateSchemaVersion is the version this build writes. A file
// below it is upgraded by internal/services/editor_state_migration; a file
// above it is rejected rather than loaded best-effort.
const CurrentEditorStateSchemaVersion = 2

type EditorState struct {
	TemplateIdentity
	MapSettings
	PlayerSettings
	NeutralZoneSettings
	CastleSettings
	GenerationSettings
	GameRuleSettings
	ContentSettings
	ManualEditSettings
	SchemaOptions
}
