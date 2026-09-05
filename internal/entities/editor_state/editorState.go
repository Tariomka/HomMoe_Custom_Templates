package editor_state

const CurrentEditorStateSchemaVersion = 1

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
