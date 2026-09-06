package editor_state_v1

// CastleSettings is the frozen v1 snapshot - see editorState.go. Never edit.
type CastleSettings struct {
	AdvancedMode                bool `json:"advancedMode"`
	PlayerOwnedCastles          int  `json:"playerOwnedCastles"`
	PlayerZoneCastles           int  `json:"playerCastles"`
	NeutralZoneCastles          int  `json:"neutralCastles"`
	HubZoneCastles              int  `json:"hubCastles"`
	NeutralLowestCastlesPerZone int  `json:"neutralLowestCastlesPerZone"`
	NeutralLowCastlesPerZone    int  `json:"neutralLowCastlesPerZone"`
	NeutralMediumCastlesPerZone int  `json:"neutralMedCastlesPerZone"`
	NeutralHighCastlesPerZone   int  `json:"neutralHighCastlesPerZone"`
	MatchPlayerCastleFactions   bool `json:"matchPlayerCastleFactions"`
}
