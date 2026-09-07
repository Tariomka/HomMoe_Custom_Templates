// Package editor_state_v1 is a frozen snapshot of the editor-state schema as it
// stood at schema version 1. Nothing here may ever be edited.
package editor_state_v1

const SchemaVersion = 1 // Version of this snapshot.

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
