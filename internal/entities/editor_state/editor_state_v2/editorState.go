// Package editor_state_v2 is a frozen snapshot of the editor-state schema as it
// stood at schema version 2. Nothing here may ever be edited.
package editor_state_v2

const SchemaVersion = 2 // Version of this snapshot.

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
