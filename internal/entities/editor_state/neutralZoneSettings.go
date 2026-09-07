package editor_state

// NeutralZoneSettings holds how many neutral zones exist, how they split by
// quality tier and whether abandoned outposts are spawned.
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
