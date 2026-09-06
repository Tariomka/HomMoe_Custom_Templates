// Package editor_state_v1 is a frozen snapshot of the editor-state schema as it
// stood at schema version 1. Nothing here may ever be edited: it exists so that
// a later change to a current editor_state struct can never silently redefine
// what a v1 .gen.json meant. It is read by
// internal/services/editor_state_migration, which is where MigrateToV2 lives.
//
// The snapshot deliberately stops at the editor-state boundary. ManualZoneSave
// and ManualConnectionSave still name the live entities.Zone and
// entities.Connection because those are the game's .rmg.json vocabulary and
// version with the game, not with this file format.
package editor_state_v1

// SchemaVersion is the version this snapshot describes.
const SchemaVersion = 1

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
