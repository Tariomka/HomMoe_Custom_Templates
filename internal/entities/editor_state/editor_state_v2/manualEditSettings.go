package editor_state_v2

// ManualEditSettings holds the authoritative snapshot of the zones and
// connections edited in the manual zone editor.
type ManualEditSettings struct {
	ManualZones       []ManualZoneSave       `json:"manualZones,omitempty"`
	ManualConnections []ManualConnectionSave `json:"manualConnections,omitempty"`
}
