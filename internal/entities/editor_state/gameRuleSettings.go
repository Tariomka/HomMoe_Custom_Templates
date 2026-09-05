package editor_state

// GameRuleSettings holds the victory condition and the optional scenario rules
// layered on top of it.
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
