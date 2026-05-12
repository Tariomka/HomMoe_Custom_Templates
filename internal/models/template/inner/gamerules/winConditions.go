package gamerules

// WinConditions enumerates every observed win-condition toggle and tuning value.
// All fields are optional in the source JSON; absent fields keep their Go zero value.
type WinConditions struct {
	Classic         bool `json:"classic,omitempty"`
	Desertion       bool `json:"desertion,omitempty"`
	DesertionDay    int  `json:"desertionDay,omitempty"`
	DesertionValue  int  `json:"desertionValue,omitempty"`
	HeroLighting    bool `json:"heroLighting,omitempty"`
	HeroLightingDay int  `json:"heroLightingDay,omitempty"`

	LostStartCity    bool `json:"lostStartCity,omitempty"`
	LostStartCityDay int  `json:"lostStartCityDay,omitempty"`
	LostStartHero    bool `json:"lostStartHero,omitempty"`

	CityHold     bool `json:"cityHold,omitempty"`
	CityHoldDays int  `json:"cityHoldDays,omitempty"`

	GladiatorArena                       bool   `json:"gladiatorArena,omitempty"`
	GladiatorArenaRegistrationStartWork  bool   `json:"gladiatorArenaRegistrationStartWork,omitempty"`
	GladiatorArenaRegistrationStartFight bool   `json:"gladiatorArenaRegistrationStartFight,omitempty"`
	GladiatorArenaDaysDelayStart         int    `json:"gladiatorArenaDaysDelayStart,omitempty"`
	GladiatorArenaCountDay               int    `json:"gladiatorArenaCountDay,omitempty"`
	ChampionSelectRule                   string `json:"championSelectRule,omitempty"`

	// Tournament-mode ("Chosen One", "Exodus", "Massacre", "Sprint") configuration.
	Tournament             bool  `json:"tournament,omitempty"`
	TournamentPointsToWin  int   `json:"tournamentPointsToWin,omitempty"`
	TournamentSaveArmy     bool  `json:"tournamentSaveArmy,omitempty"`
	TournamentDays         []int `json:"tournamentDays,omitempty"`
	TournamentAnnounceDays []int `json:"tournamentAnnounceDays,omitempty"`
}
