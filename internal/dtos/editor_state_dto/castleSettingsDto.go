package editor_state_dto

type CastleSettingsDto struct {
	AdvancedMode                bool
	PlayerOwnedCastles          int
	PlayerZoneCastles           int
	NeutralZoneCastles          int
	HubZoneCastles              int
	NeutralLowestCastlesPerZone int
	NeutralLowCastlesPerZone    int
	NeutralMediumCastlesPerZone int
	NeutralHighCastlesPerZone   int
	MatchPlayerCastleFactions   bool
}
