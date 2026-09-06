package editor_state_v1

// NeutralZoneSettings is the frozen v1 snapshot - see editorState.go. Never edit.
type NeutralZoneSettings struct {
	NeutralZoneCount           int  `json:"neutralZoneCount"`
	SpawnAbandonedOutposts     bool `json:"spawnAbandonedOutposts"`
	AbandonedOutpostCount      int  `json:"abandonedOutpostCount"`
	NeutralLowestNoCastleCount int  `json:"neutralLowestNoCastle"`
	NeutralLowestCastleCount   int  `json:"neutralLowestCastle"`
	NeutralLowNoCastleCount    int  `json:"neutralLowNoCastle"`
	NeutralLowCastleCount      int  `json:"neutralLowCastle"`
	NeutralMediumNoCastleCount int  `json:"neutralMediumNoCastle"`
	NeutralMediumCastleCount   int  `json:"neutralMediumCastle"`
	NeutralHighNoCastleCount   int  `json:"neutralHighNoCastle"`
	NeutralHighCastleCount     int  `json:"neutralHighCastle"`
}
