package editor_state_v1

// PlayerSettings is the frozen v1 snapshot - see editorState.go. Never edit.
type PlayerSettings struct {
	PlayerCount        int `json:"playerCount"`
	HeroCountMin       int `json:"heroMin"`
	HeroCountMax       int `json:"heroMax"`
	HeroCountIncrement int `json:"heroIncrement"`
}
