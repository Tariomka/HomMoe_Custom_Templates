package editor_state_v2

// PlayerSettings holds the player count and the per-player hero limits.
type PlayerSettings struct {
	PlayerCount        int `json:"playerCount"`
	HeroCountMin       int `json:"heroMin"`
	HeroCountMax       int `json:"heroMax"`
	HeroCountIncrement int `json:"heroIncrement"`
}
