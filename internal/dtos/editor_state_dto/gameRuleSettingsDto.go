package editor_state_dto

type GameRuleSettingsDto struct {
	VictoryCondition             string
	FactionLawXpPercent          int
	AstrologyXpPercent           int
	LostStartCity                bool
	LostStartCityDay             int
	LostStartHero                bool
	CityHold                     bool
	CityHoldDays                 int
	GladiatorArena               bool
	GladiatorArenaDaysDelayStart int
	GladiatorArenaCountDay       int
	Tournament                   bool
	TournamentFirstTournamentDay int
	TournamentInterval           int
	TournamentPointsToWin        int
	TournamentSaveArmy           bool
}
