package editor_state

// CastleSettings holds every castle-count option. AdvancedMode belongs here
// because it gates which of the neutral options below are read.
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
