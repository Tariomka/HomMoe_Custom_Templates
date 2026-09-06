package editor_state_v1

// GameRuleSettings is the frozen v1 snapshot - see editorState.go. Never edit.
type GameRuleSettings struct {
	VictoryCondition             string `json:"victoryCondition"`
	FactionLawXpPercent          int    `json:"factionLawsExp"`
	AstrologyXpPercent           int    `json:"astrologyExp"`
	LostStartCity                bool   `json:"lostStartCity"`
	LostStartCityDay             int    `json:"lostStartCityDay"`
	LostStartHero                bool   `json:"lostStartHero"`
	CityHold                     bool   `json:"cityHold"`
	CityHoldDays                 int    `json:"cityHoldDays"`
	GladiatorArena               bool   `json:"gladiatorArena"`
	GladiatorArenaDaysDelayStart int    `json:"gladiatorArenaDaysDelayStart"`
	GladiatorArenaCountDay       int    `json:"gladiatorArenaCountDay"`
	Tournament                   bool   `json:"tournament"`
	TournamentFirstTournamentDay int    `json:"tournamentFirstTournamentDay"`
	TournamentInterval           int    `json:"tournamentInterval"`
	TournamentPointsToWin        int    `json:"tournamentPointsToWin"`
	TournamentSaveArmy           bool   `json:"tournamentSaveArmy"`
}
