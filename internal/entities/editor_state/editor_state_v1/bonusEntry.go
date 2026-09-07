package editor_state_v1

// BonusEntry is the frozen v1 snapshot - see editorState.go. Never edit.
type BonusEntry struct {
	PresetType     BonusPresetType `json:"presetType"`
	ReceiverFilter string          `json:"receiverFilter"`
	Param          string          `json:"param,omitempty"`
	Param2         string          `json:"param2,omitempty"`
}
