package editor_state

// CurrentEditorStateSchemaVersion is the version this build writes. A file
// below it is upgraded by internal/services/file_service/editor_state_migrator;
// a file above it is rejected rather than loaded best-effort.
const CurrentEditorStateSchemaVersion = 3

// EditorState keeps identity, map and schema fields at the root and every other
// settings group in its own section.
type EditorState struct {
	SchemaOptions
	TemplateIdentity
	MapSettings

	Players      PlayerSettings      `json:"players"`
	NeutralZones NeutralZoneSettings `json:"neutralZones"`
	Castles      CastleSettings      `json:"castles"`
	Generation   GenerationSettings  `json:"generation"`
	GameRules    GameRuleSettings    `json:"gameRules"`
	Content      ContentSettings     `json:"content"`
	ManualEdits  ManualEditSettings  `json:"manualEdits,omitzero"`
}
