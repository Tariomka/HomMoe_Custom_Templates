package template

// GameRules describes the gameplay rules for the template.
//
// A few templates (e.g. "Symmetry", "Jebus Outcast") declare gladiator-arena /
// global-bans fields as siblings of `winConditions` inside `gameRules` instead
// of nesting them properly; those fields are mirrored here for tolerant parsing.
type GameRules struct {
	HeroCountMin       int  `json:"heroCountMin"`
	HeroCountMax       int  `json:"heroCountMax"`
	HeroCountIncrement int  `json:"heroCountIncrement"`
	HeroHireBan        bool `json:"heroHireBan"`
	EncounterHoles     bool `json:"encounterHoles"`
	TournamentRules    bool `json:"tournamentRules,omitempty"`

	Bonuses       BonusList     `json:"bonuses,omitempty"`
	WinConditions WinConditions `json:"winConditions"`

	// Mirror of gladiator-arena fields when the template author placed them as
	// siblings of `winConditions` rather than inside it.
	GladiatorArena                       bool   `json:"gladiatorArena,omitempty"`
	GladiatorArenaRegistrationStartWork  bool   `json:"gladiatorArenaRegistrationStartWork,omitempty"`
	GladiatorArenaRegistrationStartFight bool   `json:"gladiatorArenaRegistrationStartFight,omitempty"`
	GladiatorArenaDaysDelayStart         int    `json:"gladiatorArenaDaysDelayStart,omitempty"`
	GladiatorArenaCountDay               int    `json:"gladiatorArenaCountDay,omitempty"`
	ChampionSelectRule                   string `json:"championSelectRule,omitempty"`

	// Some templates declare `globalBans` inside `gameRules` instead of at root.
	GlobalBans *GlobalBans `json:"globalBans,omitempty"`

	FactionLawsExpModifier float64 `json:"factionLawsExpModifier,omitempty"`
	AstrologyExpModifier   float64 `json:"astrologyExpModifier,omitempty"`
}

// Bonus describes a starting bonus granted to a side / hero.
// `parameters` contents are type-dependent; all observed values are JSON strings.
type Bonus struct {
	SID            string   `json:"sid"`
	ReceiverSide   int      `json:"receiverSide"`
	ReceiverFilter string   `json:"receiverFilter,omitempty"`
	Parameters     []string `json:"parameters"`
}

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
