package editor_state_v1

// ManualEditSettings is the frozen v1 snapshot - see editorState.go. Never edit.
type ManualEditSettings struct {
	ManualZones       []ManualZoneSave       `json:"manualZones,omitempty"`
	ManualConnections []ManualConnectionSave `json:"manualConnections,omitempty"`
}
